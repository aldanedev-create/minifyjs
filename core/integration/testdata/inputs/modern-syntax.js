const arrow = (x) => x * 2;
const { a, b } = { a: 1, b: 2 };
const [c, ...rest] = [1, 2, 3, 4];
const template = `hello ${arrow(a)}`;
const optional = c?.toString?.();
const nullish = optional ?? "default";

class Foo {
    #private = 1;
    static counter = 0;
    constructor(x) {
        this.x = x;
        Foo.counter++;
    }
    get value() { return this.#private; }
}

async function fetchData() {
    return await Promise.resolve(42);
}

function* gen() {
    yield 1;
    yield 2;
}

for (const item of [1, 2, 3]) {
    console.log(item);
}

let obj = { a: 1 };
obj.b ||= 2;
obj.c &&= 3;
obj.d ??= 4;