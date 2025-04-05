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

		want := "if (true) {\n  console.log('hello')\n}\n"
		assert.Equal(t, want, g.String())
	})

	t.Run("CustomSpaces", func(t *testing.T) {
		g := NewGenKit().WithSpaces(4)
		g.Line("if (true) {").
			Indent().
			Line("console.log('hello')").
			Dedent().
			Line("}")

		want := "if (true) {\n    console.log('hello')\n}\n"
		assert.Equal(t, want, g.String())
	})

	t.Run("WithTabs", func(t *testing.T) {
		g := NewGenKit().WithTabs()
		g.Line("if (true) {").
			Indent().
			Line("console.log('hello')").
			Dedent().
			Line("}")

		want := "if (true) {\n\tconsole.log('hello')\n}\n"
		assert.Equal(t, want, g.String())
	})

	t.Run("Empty lines", func(t *testing.T) {
		g := NewGenKit().WithSpaces(2)
		g.Line("// This is a comment").
			Break().
			Break().
			Line("// This is other comment")

		want := "// This is a comment\n\n\n// This is other comment\n"
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

		want := "if (true) {\n\tconsole.log('hello')\n}\n"
		assert.Equal(t, want, g.String())
	})

	t.Run("Test Without Chain of Methods", func(t *testing.T) {
		g := NewGenKit().WithSpaces(2).WithTabs()
		g.Line("if (true) {")
		g.Indent()
		g.Line("console.log('hello')")
		g.Dedent()
		g.Line("}")

		want := "if (true) {\n\tconsole.log('hello')\n}\n"
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

		want := `if (true) {
  console.log('hello')
  console.log('world')
}
`
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

		want := `function example() {
  if (condition) {
    console.log('condition true')
  } else {
    console.log('condition false')
  }
}
`
		assert.Equal(t, want, g.String())
	})
}

func TestLinef(t *testing.T) {
	t.Run("SimpleFormat", func(t *testing.T) {
		g := NewGenKit().WithSpaces(2)
		g.Linef("const greeting = %q", "Hello, World!")

		want := `const greeting = "Hello, World!"
`
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

		want := `interface User {
  id: string
  age: number
}
`
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

		want := `type User struct {
  ID   int
  Name string
}
`
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

		want := `class User:
    def __init__(self, name, age):
        self.name = name
        self.age = age

    def greet(self):
        return f"Hello, {self.name}!"
`
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

		want := `class User
  def initialize(name, age)
    @name = name
    @age = age
  end

  def greet
    puts "Hello, #{@name}!"
  end
end
`
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

		want := `public class User {
    private String name;
    private int age;

    public User(String name, int age) {
        this.name = name;
        this.age = age;
    }

    public String greet() {
        return "Hello, " + name + "!";
    }
}
`
		assert.Equal(t, want, g.String())
	})
}

func TestInline(t *testing.T) {
	t.Run("SimpleInline", func(t *testing.T) {
		g := NewGenKit().WithSpaces(2)
		g.Inline("console.log('hello')")

		want := "console.log('hello')"
		assert.Equal(t, want, g.String())
	})

	t.Run("InlineWithIndentation", func(t *testing.T) {
		g := NewGenKit().WithSpaces(2)
		g.Indent().Inline("console.log('hello')")

		want := "  console.log('hello')"
		assert.Equal(t, want, g.String())
	})

	t.Run("InlineWithMultipleLines", func(t *testing.T) {
		g := NewGenKit().WithSpaces(2)
		g.Inline("console.log('hello')").Inline("console.log('world')")

		want := "console.log('hello')console.log('world')"
		assert.Equal(t, want, g.String())
	})
}

func TestInlinef(t *testing.T) {
	t.Run("SimpleInlinef", func(t *testing.T) {
		g := NewGenKit().WithSpaces(2)
		g.Inlinef("const greeting = %q", "Hello, World!")

		want := "const greeting = \"Hello, World!\""
		assert.Equal(t, want, g.String())
	})

	t.Run("InlinefWithIndentation", func(t *testing.T) {
		g := NewGenKit().WithSpaces(2)
		g.Indent().Inlinef("const greeting = %q", "Hello, World!")

		want := "  const greeting = \"Hello, World!\""
		assert.Equal(t, want, g.String())
	})

	t.Run("InlinefWithMultipleLines", func(t *testing.T) {
		g := NewGenKit().WithSpaces(2)
		g.Inlinef("const greeting = %q", "Hello, World!").Inlinef("console.log(%q)", "Hello, World!")

		want := "const greeting = \"Hello, World!\"console.log(\"Hello, World!\")"
		assert.Equal(t, want, g.String())
	})
}

func TestSublines(t *testing.T) {
	t.Run("If a line contains newlines, each line will be properly indented", func(t *testing.T) {
		g := NewGenKit().WithSpaces(2)
		g.Line("function test() {")
		g.Block(func() {
			g.Line("console.log('hello')\nconsole.log('world')")
		})
		g.Line("}")

		want := `function test() {
  console.log('hello')
  console.log('world')
}
`
		assert.Equal(t, want, g.String())
	})
}

func TestInlineContinuation(t *testing.T) {
	t.Run("Inline multiple times on same line", func(t *testing.T) {
		g := NewGenKit().WithSpaces(2)
		g.Line("const x = ").
			Indent().
			Inline("1").
			Inline(" + ").
			Inline("2;")

		want := "const x = \n  1 + 2;"
		assert.Equal(t, want, g.String())
	})

	t.Run("Inline after Line", func(t *testing.T) {
		g := NewGenKit().WithSpaces(2)
		g.Line("line1")
		g.Inline("line2")

		want := "line1\nline2"
		assert.Equal(t, want, g.String())
	})

	t.Run("Line after Inline", func(t *testing.T) {
		g := NewGenKit().WithSpaces(2)
		g.Inline("line1")
		g.Line("line2")

		want := "line1line2\n"
		assert.Equal(t, want, g.String())
	})
}
