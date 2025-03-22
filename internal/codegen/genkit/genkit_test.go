package genkit

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBasicIndentation(t *testing.T) {
	t.Run("DefaultIndentation", func(t *testing.T) {
		g := NewGenKit()
		g.Line("if (true) {").
			Indent().
			Line("console.log('hello')").
			Dedent().
			Line("}")

		want := "\nif (true) {\n  console.log('hello')\n}"
		assert.Equal(t, want, g.String())
	})

	t.Run("CustomSpaces", func(t *testing.T) {
		g := NewGenKit().WithSpaces(4)
		g.Line("if (true) {").
			Indent().
			Line("console.log('hello')").
			Dedent().
			Line("}")

		want := "\nif (true) {\n    console.log('hello')\n}"
		assert.Equal(t, want, g.String())
	})

	t.Run("WithTabs", func(t *testing.T) {
		g := NewGenKit().WithTabs()
		g.Line("if (true) {").
			Indent().
			Line("console.log('hello')").
			Dedent().
			Line("}")

		want := "\nif (true) {\n\tconsole.log('hello')\n}"
		assert.Equal(t, want, g.String())
	})

	t.Run("Empty lines", func(t *testing.T) {
		g := NewGenKit().WithSpaces(2)
		g.Line("// This is a comment").
			Line().
			Line().
			Line("// This is other comment")

		want := "\n// This is a comment\n\n\n// This is other comment"
		assert.Equal(t, want, g.String())
	})
}

func TestChainable(t *testing.T) {
	t.Run("Test With Chain of Methods", func(t *testing.T) {
		g := NewGenKit().WithSpaces(2).WithTabs()
		g.Line("if (true) {").
			Indent().
			Line("console.log('hello')").
			Dedent().
			Line("}")

		want := "\nif (true) {\n\tconsole.log('hello')\n}"
		assert.Equal(t, want, g.String())
	})

	t.Run("Test Without Chain of Methods", func(t *testing.T) {
		g := NewGenKit().WithSpaces(2).WithTabs()
		g.Line("if (true) {")
		g.Indent()
		g.Line("console.log('hello')")
		g.Dedent()
		g.Line("}")

		want := "\nif (true) {\n\tconsole.log('hello')\n}"
		assert.Equal(t, want, g.String())
	})
}

func TestBlock(t *testing.T) {
	t.Run("SimpleBlock", func(t *testing.T) {
		g := NewGenKit().WithSpaces(2)
		g.Line("if (true) {").
			Block(func() {
				g.Line("console.log('hello')").
					Line("console.log('world')")
			}).
			Line("}")

		want := `
if (true) {
  console.log('hello')
  console.log('world')
}`
		assert.Equal(t, want, g.String())
	})

	t.Run("NestedBlocks", func(t *testing.T) {
		g := NewGenKit().WithSpaces(2)
		g.Line("function example() {").
			Block(func() {
				g.Line("if (condition) {").
					Block(func() {
						g.Line("console.log('condition true')")
					}).
					Line("} else {").
					Block(func() {
						g.Line("console.log('condition false')")
					}).
					Line("}")
			}).
			Line("}")

		want := `
function example() {
  if (condition) {
    console.log('condition true')
  } else {
    console.log('condition false')
  }
}`
		assert.Equal(t, want, g.String())
	})
}

func TestLinef(t *testing.T) {
	t.Run("SimpleFormat", func(t *testing.T) {
		g := NewGenKit().WithSpaces(2)
		g.Linef("const greeting = %q", "Hello, World!")

		want := `
const greeting = "Hello, World!"`
		assert.Equal(t, want, g.String())
	})

	t.Run("TypeScriptInterface", func(t *testing.T) {
		g := NewGenKit().WithSpaces(2)
		g.Linef("interface %s {", "User").
			Indent().
			Linef("id: %s", "string").
			Linef("age: %s", "number").
			Dedent().
			Line("}")

		want := `
interface User {
  id: string
  age: number
}`
		assert.Equal(t, want, g.String())
	})

	t.Run("GoStruct", func(t *testing.T) {
		g := NewGenKit().WithSpaces(2)
		g.Linef("type %s struct {", "User").
			Indent().
			Linef("ID   %s", "int").
			Linef("Name %s", "string").
			Dedent().
			Line("}")

		want := `
type User struct {
  ID   int
  Name string
}`
		assert.Equal(t, want, g.String())
	})
}

func TestMultiLanguageExamples(t *testing.T) {
	t.Run("PythonClass", func(t *testing.T) {
		g := NewGenKit().WithSpaces(4)
		g.Line("class User:").
			Indent().
			Line("def __init__(self, name, age):").
			Indent().
			Line("self.name = name").
			Line("self.age = age").
			Dedent().
			Line("").
			Line("def greet(self):").
			Indent().
			Line("return f\"Hello, {self.name}!\"").
			Dedent().
			Dedent()

		want := `
class User:
    def __init__(self, name, age):
        self.name = name
        self.age = age

    def greet(self):
        return f"Hello, {self.name}!"`
		assert.Equal(t, want, g.String())
	})

	t.Run("RubyClass", func(t *testing.T) {
		g := NewGenKit().WithSpaces(2)
		g.Line("class User").
			Indent().
			Line("def initialize(name, age)").
			Indent().
			Line("@name = name").
			Line("@age = age").
			Dedent().
			Line("end").
			Line("").
			Line("def greet").
			Indent().
			Line("puts \"Hello, #{@name}!\"").
			Dedent().
			Line("end").
			Dedent().
			Line("end")

		want := `
class User
  def initialize(name, age)
    @name = name
    @age = age
  end

  def greet
    puts "Hello, #{@name}!"
  end
end`
		assert.Equal(t, want, g.String())
	})

	t.Run("JavaClass", func(t *testing.T) {
		g := NewGenKit().WithSpaces(4)
		g.Linef("public class %s {", "User").
			Indent().
			Line("private String name;").
			Line("private int age;").
			Line("").
			Linef("public %s(String name, int age) {", "User").
			Indent().
			Line("this.name = name;").
			Line("this.age = age;").
			Dedent().
			Line("}").
			Line("").
			Line("public String greet() {").
			Indent().
			Line("return \"Hello, \" + name + \"!\";").
			Dedent().
			Line("}").
			Dedent().
			Line("}")

		want := `
public class User {
    private String name;
    private int age;

    public User(String name, int age) {
        this.name = name;
        this.age = age;
    }

    public String greet() {
        return "Hello, " + name + "!";
    }
}`
		assert.Equal(t, want, g.String())
	})
}
