// This file exists so verify_no_node.sh has something to minify.
// The script proves MinifyJS works on a machine with no Node.js.

function add(a, b) {
    return a + b;
}

function greet(name) {
    return "hello, " + name;
}

console.log(greet("world"));
console.log(add(2, 3));