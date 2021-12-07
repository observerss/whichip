#! /usr/bin/env node
const { discover, listen, logger, PORT } = require('..')
const { Command } = require('commander')

function main () {
  const program = new Command()
  program.version('0.1.0')
  program
    .command('listen')
    .option('--debug', 'print debug message')
    .option('--port <PORT>', 'port to listen on', parsePort, PORT)
    .action(listen)

  program
    .command('discover')
    .option('--debug', 'print debug message', false)
    .option('--port <PORT>', 'remote port', parsePort, PORT)
    .option('--timeout <TIMEOUT>', 'network timeout', parseTimeout, 1.0)
    .option('--all', 'wait all responses until timeout (default to return on first response)')
    .action((opts, prog) => {
      opts.log = true
      discover(opts, prog).catch(console.error)})

  program.parse(process.argv)
}

function parsePort (value, dummyPrevious) {
  const parsedValue = parseInt(value, 10)
  if (isNaN(parsedValue)) {
    logger.fatal('`port` should be a number')
    process.exit(1)
  }
  return parsedValue
}

function parseTimeout (value, dummyPrevious) {
  const parsedValue = parseFloat(value)
  if (isNaN(parsedValue)) {
    logger.fatal('`timeout` should be a number')
    process.exit(1)
  }
  return parsedValue
}

main()