package minify_test

import (
	"testing"

	"github.com/minifyjs/minifyjs/core/integration/helpers"
)

func TestArrowFunctions(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{
		Stdin: []byte("const f = (x) => x * 2;"),
	})
	helpers.AssertExitCode(t, r, 0)
	helpers.AssertContains(t, string(r.Stdout), "=>")
}

func TestTemplateLiterals(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{
		Stdin: []byte("const s = `hello ${name}`;"),
	})
	helpers.AssertExitCode(t, r, 0)
	helpers.AssertContains(t, string(r.Stdout), "`")
}

func TestDestructuringObject(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{
		Stdin: []byte("const { a, b } = obj;"),
	})
	helpers.AssertExitCode(t, r, 0)
}

func TestDestructuringArray(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{
		Stdin: []byte("const [x, y] = arr;"),
	})
	helpers.AssertExitCode(t, r, 0)
}

func TestSpreadOperator(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{
		Stdin: []byte("const c = [...a, ...b];"),
	})
	helpers.AssertExitCode(t, r, 0)
}

func TestRestParameters(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{
		Stdin: []byte("function f(...args) { return args.length; }"),
	})
	helpers.AssertExitCode(t, r, 0)
}

func TestDefaultParameters(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{
		Stdin: []byte("function f(x = 1) { return x; }"),
	})
	helpers.AssertExitCode(t, r, 0)
}

func TestClasses(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{
		Stdin: []byte("class A { constructor(x) { this.x = x; } }"),
	})
	helpers.AssertExitCode(t, r, 0)
	helpers.AssertContains(t, string(r.Stdout), "class")
}

func TestClassInheritance(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{
		Stdin: []byte("class B extends A { constructor() { super(); } }"),
	})
	helpers.AssertExitCode(t, r, 0)
}

func TestPrivateFields(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{
		Stdin: []byte("class A { #x = 1; getX() { return this.#x; } }"),
	})
	helpers.AssertExitCode(t, r, 0)
	helpers.AssertContains(t, string(r.Stdout), "#x")
}

func TestStaticFields(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{
		Stdin: []byte("class A { static x = 1; }"),
	})
	helpers.AssertExitCode(t, r, 0)
}

func TestGettersSetters(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{
		Stdin: []byte("class A { get x() { return 1; } set x(v) {} }"),
	})
	helpers.AssertExitCode(t, r, 0)
}

func TestAsyncFunction(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{
		Stdin: []byte("async function f() { return 1; }"),
	})
	helpers.AssertExitCode(t, r, 0)
	helpers.AssertContains(t, string(r.Stdout), "async")
}

func TestAwaitExpression(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{
		Stdin: []byte("async function f() { await g(); }"),
	})
	helpers.AssertExitCode(t, r, 0)
}

func TestGenerators(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{
		Stdin: []byte("function* gen() { yield 1; }"),
	})
	helpers.AssertExitCode(t, r, 0)
}

func TestForOfLoop(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{
		Stdin: []byte("for (const x of arr) { console.log(x); }"),
	})
	helpers.AssertExitCode(t, r, 0)
}

func TestOptionalChaining(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{
		Stdin: []byte("const x = a?.b?.c;"),
	})
	helpers.AssertExitCode(t, r, 0)
}

func TestNullishCoalescing(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{
		Stdin: []byte("const x = a ?? b;"),
	})
	helpers.AssertExitCode(t, r, 0)
}

func TestLogicalAssignment(t *testing.T) {
	r := helpers.Run(t, helpers.RunOptions{
		Stdin: []byte("a ||= b; a &&= c; a ??= d;"),
	})
	helpers.AssertExitCode(t, r, 0)
}