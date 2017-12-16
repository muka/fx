const ascii = require('ascii-art')

console.log(JSON.stringify(process.argv))

const hello = process.argv[2] || 'fx'

//NOTE this is used in unit test to validate input / output
//     (see ../../docker-api/exec_test.go)
const text = `hello ${hello}!`
ascii.font(text, 'Doom', (out) => {
    console.log(text)
    console.log(out)
})
