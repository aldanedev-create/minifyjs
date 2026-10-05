// A small module with the common shapes found in application code.

function add(a, b) {
    return a + b;
}

function subtract(a, b) {
    return a - b;
}

function multiply(a, b) {
    return a * b;
}

function divide(a, b) {
    if (b === 0) {
        throw new Error("division by zero");
    }
    return a / b;
}

const operations = {
    add: add,
    subtract: subtract,
    multiply: multiply,
    divide: divide,
};

class Calculator {
    constructor() {
        this.history = [];
    }

    compute(op, a, b) {
        const fn = operations[op];
        if (!fn) {
            throw new Error("unknown operation: " + op);
        }
        const result = fn(a, b);
        this.history.push({ op: op, a: a, b: b, result: result });
        return result;
    }

    getHistory() {
        return this.history.slice();
    }
}

export { add, subtract, multiply, divide, Calculator };