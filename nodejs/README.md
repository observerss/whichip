# NodeJS Client

## Install

```bash
npm install -g whichip
```

## Usage

### In Terminal

```bash
$ whichip discover --timeout 0.2 --all --debug
#10.86.8.222
# or
#not found
```

### In Code


```bash
npm install whichip
```

Then,

```javascript
const { discover } = require('whichip')

opts = {
  debug: false,  // print debug message
  port: 53535,   // remote server port
  all: true,    // true => wait until timeout; false => return on first ip
  timeout: 0.2,  // discover timeout, seconds
}
discover(opts)
  .then((ips) => {console.log(ips)})
  .catch(console.error)
// [ '10.86.8.222' ]
// -- or --
// not found
```

