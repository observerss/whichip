const dgram = require('dgram')
const logger = require('pino')()

const PORT = 53535
const CLIENT_MSG = 'pfg_ip_broadcast_cl'
const SERVER_MSG = 'pfg_ip_response_serv'
let clientsToClose = []
let timeouts = []

/*
  The listen command
  params
    opts: {
      debug: false, // print debug message
      port: 53535, // the port to listen on
    }
  returns
    null
 */
function listen (opts, program) {
  opts = Object.assign({ debug: false, port: PORT }, opts || {})
  setupLogger(opts.debug)

  const server = dgram.createSocket('udp4')
  server.on('error', (err) => {
    logger.debug(`server error:\n${err.stack}`)
    server.close()
  })

  server.on('message', (msg, rinfo) => {
    logger.debug(`server got: ${msg} from ${rinfo.address}:${rinfo.port}`)
    if (msg.toString() != CLIENT_MSG) {
      logger.debug('unknown message')
    } else {
      server.send(SERVER_MSG, rinfo.port, rinfo.address)
      logger.debug(`sent back ${SERVER_MSG} to ${rinfo.address}:${rinfo.port}`)
    }
  })

  server.on('listening', () => {
    const address = server.address()
    logger.debug(`server listening ${address.address}:${address.port}`)
  })

  server.bind(opts.port)
}

/*
  The discover command
  params
    opts: {
      debug: false,  // print debug message
      port: 53535,   // remote server port
      all: true,    // true => wait until timeout; false => return on first ip
      timeout: 0.5,  // discover timeout
      log: false,    // print ip to console
    }
  returns
    Promise<string[]>, a list of ips, e.g. Promise<['10.86.8.222']>
 */
function discover (opts, program) {
  opts = Object.assign({ debug: false, port: PORT, all: true, timeout: 0.2 }, opts || {})
  setupLogger(opts.debug)

  return new Promise((resolve, reject) => {
    let jobs = []
    for (let inet of usableAddresses()) {
      jobs.push(discoverOnAddress(inet, opts.port, opts.timeout, opts.all))
    }

    if (opts.all) {
      Promise.allSettled(jobs)
        .then((values) => {
          let ips = values
            .filter(x => x.status === 'fulfilled')
            .map(x => x.value)
            .flat()
          clearAll()
          if (ips.length > 0) {
            for (let ip of ips) {
              if (opts.log)
                console.log(ip)
            }
            resolve(ips)
          } else {
            reject('not found')
          }
        })
    } else {
      promiseAny(jobs).then((ips) => {
        if (opts.log)
          console.log(ips[0])
        resolve(ips)
      }).catch((error) => {
          clearAll()
          reject('not found')
        }
      )
    }
  })
}

/*
  Find all usable IPs

  returns
    string[]: usable IPs, e.g. ['10.86.8.99', '172.17.48.1']
 */
function usableAddresses () {
  return Object.values(require('os').networkInterfaces())
    .flat()
    .filter(x => x.family === 'IPv4')
    .filter(x => x.internal === false)
    .map(x => x.address)
}

/*
  Discover ip(s) on specified address

  params
    inet: the ipv4 local addr, e.g. 10.86.2.99
    port: the remote port, e.g. 53535
    timeout: discover timeout, e.g. 1.0
    all: true => wait until timeout; false => return on first ip
  returns
    Promise<string[]>: a list of ips, e.g. Promise<['10.86.8.222']>
 */
function discoverOnAddress (inet, port, timeout, all) {
  let ips = []
  let onResolve = null
  let broadcastAddr = '255.255.255.255'
  let closed = false

  const client = dgram.createSocket('udp4')
  clientsToClose.push(client)
  client.on('error', (err) => {
    logger.debug(`client error:\n${err.stack}`)
    client.close()
  })

  client.on('message', (msg, rinfo) => {
    logger.debug(`client got: ${msg} from ${rinfo.address}:${rinfo.port}`)
    ips.push(rinfo.address)
    if (!all) {
      clearAll()
      onResolve(ips)
    }
  })

  client.on('listening', () => {
    client.setBroadcast(true)
    client.send(CLIENT_MSG, port, broadcastAddr, (error, msg) => {
      if (error) {
        logger.debug(`client error: cannot send broadcast to ${broadcastAddr}:${port} on address ${inet}, ${msg}`)
      } else {
        logger.debug(`client sent: ${CLIENT_MSG} to ${broadcastAddr}:${port} on address ${inet}`)
      }
    })
  })

  client.bind(0, inet)

  return new Promise((resolve, reject) => {
    onResolve = resolve
    timeouts.push(setTimeout(() => {
      if (ips.length > 0) {
        resolve(ips)
      } else {
        reject('timed out')
      }
    }, timeout * 1000))
  })
}

function setupLogger (debug) {
  if (debug) {
    logger.level = 'debug'
  } else {
    logger.level = 'info'
  }
}

function clearAll () {
  for (let cli of clientsToClose) {
    try {
      cli.close()
    } catch {}
  }
  for (let timeout of timeouts) {
    clearTimeout(timeout)
  }
}

const promiseAny = async function (iterable) {
  return Promise.all(
    [...iterable].map(promise => {
      return new Promise((resolve, reject) =>
        Promise.resolve(promise).then(reject, resolve)
      )
    })
  ).then(
    errors => Promise.reject(errors),
    value => Promise.resolve(value)
  )
}

module.exports = {
  discover,
  listen,
  logger,
  PORT
}