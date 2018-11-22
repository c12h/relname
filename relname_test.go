package relname_test

import (
	"github.com/c12h/relname"
	"testing"
)

/*=========================== Testing CleanString ============================*/

func TestCleanString(t *testing.T) {
	testCleanString(t, "a     b\nc   ", "a b c")
	testCleanString(t, "", "")
	testCleanString(t, "\n    \r", "")
	testCleanString(t, "«»", "«»")
	testCleanString(t, "  x", "x")
	testCleanString(t, "x  ", "x")
	testCleanString(t, "  x   ", "x")
}

func testCleanString(t *testing.T, arg, expected string) {
	if actual := relname.CleanString(arg); actual != expected {
		t.Errorf("expected %-16q but got %-16q from CleanString(%q)\n",
			expected, actual, arg)
	}
}

/*=========================== Testing Name objects ===========================*/

// Test the constructors for Name, without incurring any errors.
func TestNewName123(t *testing.T) {
	expectName(t, relname.Name{}, "", "", "", "", "", 0)
	//
	Baen, err := relname.NewName1("Baen")
	check(t, err, `NewName1("Baen")`)
	expectName(t, Baen, "Baen", "Baen", "Baen", "", "", 1)
	//
	Baen_Books, err := relname.NewName1("Baen Books")
	check(t, err, `NewName1("Baen Books")`)
	expectName(t, Baen_Books, "Baen Books", "Baen Books",
		"Baen Books", "", "", 1)
	//
	Dave_F, err := relname.NewName2("Dave", "Freer")
	check(t, err, `NewName2("Dave", "Freer")`)
	expectName(t, Dave_F, "Dave Freer", "Freer, Dave",
		"Freer", "Dave", "", 2)
	//
	S_VS, err := relname.NewName2("Sydney", "Van Scyoc")
	check(t, err, `NewName2("Sydney", "Van Scyoc")`)
	expectName(t, S_VS, "Sydney Van Scyoc", "Van Scyoc, Sydney",
		"Van Scyoc", "Sydney", "", 2)
	//
	W_K_J, err := relname.NewName3("William H.  ", "      Keith", "Jr.")
	check(t, err, `NewName3("William", "Keith", "Jr.")`)
	expectName(t, W_K_J, "William H. Keith Jr.", "Keith, William H. Jr.",
		"Keith", "William H.", "Jr.", 3)
}

func check(t *testing.T, e error, call string) {
	if e != nil {
		t.Fatalf("Unexpected error from %s: %#v\n", call, e)
	}
}

func expectName(t *testing.T, name relname.Name,
	common, fileAs, surname, forename, gen string, nParts int) {
	if actual := name.NumParts(); actual != nParts {
		t.Errorf("expected %#v.NumParts() to give %d, got %d\n",
			name, nParts, actual)
	}
	expectNameStr(t, name, "Common", name.String(), common)
	expectNameStr(t, name, "FileAs", name.FileAs(), fileAs)
	expectNameStr(t, name, "Forename", name.Forename(), forename)
	expectNameStr(t, name, "Surname", name.Surname(), surname)
	expectNameStr(t, name, "Generation", name.Generation(), gen)
	expectNameStr(t, name, "String", name.String(), common)
}

func expectNameStr(t *testing.T, name relname.Name, method, actual, expected string) {
	if actual != expected {
		t.Errorf("expected %-16q but got %-16q from %#v.%s()\n",
			expected, actual, name, method)
	}
}

// Test the constructors for Name when they report errors.
func TestNewName123Errors(t *testing.T) {
	var err error
	//
	epe := &relname.EmptyPartError{}
	expectEmptyPartError(t, `&EmptyPartError{}`, epe, &relname.EmptyPartError{},
		`BUG: bad EmptyPartError value relname.EmptyPartError`+
			`{NumArgs:0, Arg1:"", Arg2:"", Arg3:""}`)
	//
	_, err = relname.NewName1(" \t ")
	expectEmptyPartError(t, `NewName1(" \t ")`, err,
		&relname.EmptyPartError{NumArgs: 1, Arg1: " \t "},
		`empty or whitespace-only argument in NewName1(" \t ")`)
	//
	_, err = relname.NewName2("", "Smith")
	expectEmptyPartError(t, `newName2("", "Smith")`, err,
		&relname.EmptyPartError{NumArgs: 2, Arg1: "", Arg2: "Smith"},
		`empty or whitespace-only argument in NewName2("", "Smith")`)
	//
	_, err = relname.NewName2(" John ", "")
	expectEmptyPartError(t, `newName2(" John ", "")`, err,
		&relname.EmptyPartError{NumArgs: 2, Arg1: " John ", Arg2: ""},
		`empty or whitespace-only argument in NewName2(" John ", "")`)
	//
	_, err = relname.NewName2(" ", "\t")
	expectEmptyPartError(t, `newName2(" ", "\t")`, err,
		&relname.EmptyPartError{NumArgs: 2, Arg1: " ", Arg2: "\t"},
		`empty or whitespace-only arguments in NewName2(" ", "\t")`)
	//
	expectEmptyPartError(t, `&EmptyPartError{99, "one", "two", "three"}`,
		&relname.EmptyPartError{99, "one", "two", "three"},
		&relname.EmptyPartError{99, "one", "two", "three"},
		`empty or whitespace-only arguments in NewName99("one", "two", "three")`)
}

func expectEmptyPartError(t *testing.T, call string, actual error,
	expected *relname.EmptyPartError, expectedText string) {
	//
	if a, ok := actual.(*relname.EmptyPartError); !ok {
		t.Errorf("%s → weird error %#v\n", call, actual)
	} else if *a != *expected {
		t.Errorf("%s → wrong value\n\tgot      %#v,\n\texpected %#v\n",
			call, a, expected)
	} // else { t.Logf("%s → right type, right value\n", call) }
	//
	if actual.Error() != expectedText {
		t.Errorf("%s → error “%s”\n\texpected “%s”\n",
			call, actual, expectedText)
	}
}

/*======================= Testing RelatedName objects ========================*/

func TestRelatedName(t *testing.T) {
	_, err := relname.NewRelatedName(relname.Name{}, "aut")
	call1 := `NewRelatedName(relname.Name{}, "aut")`
	exer1 := `&BadName{"aut"}`
	if err == nil {
		t.Errorf(`%s succeeded, should have failed`, call1)
	} else if bn, ok := err.(*relname.BadName); !ok {
		t.Errorf(`%s → wierd error %#v`, call1, err)
	} else if *bn != (relname.BadName{"aut"}) {
		t.Errorf(`%s → error %#v\n\tright type but expected %s`,
			call1, err, exer1)
	}

	Sarah, _ := relname.NewName2("Sarah \t\t\t A.  ", "Hoyt   ")
	expectBadRelocator(t, "Aut")
	expectBadRelocator(t, "aut.")
	expectBadRelocator(t, " aut")
	expectBadRelocator(t, "aut ")
	expectBadRelocator(t, " aut ")
	expectBadRelocator(t, "au")
	expectBadRelocator(t, "")
	expectBadRelocator(t, "author")

	SarahAut, _ := relname.NewRelatedName(Sarah, "aut")
	expectRelName(t, "Relator",
		SarahAut.Relator(),
		"aut")
	expectRelName(t, "Common",
		SarahAut.Common(),
		"Sarah A. Hoyt")
	expectRelName(t, "FileAs",
		SarahAut.FileAs(),
		"Hoyt, Sarah A.")
	expectRelName(t, "Forename",
		SarahAut.Forename(),
		"Sarah A.")
	expectRelName(t, "Surname",
		SarahAut.Surname(),
		"Hoyt")
	expectRelName(t, "Generation",
		SarahAut.Generation(),
		"")
	expectRelName(t, "String",
		SarahAut.String(),
		"Sarah A. Hoyt (aut)")
}

var Sarah, _ = relname.NewName2("Sarah \t\t\t A.  ", "Hoyt   ")

func expectBadRelocator(t *testing.T, relCode string) {
	_, err := relname.NewRelatedName(Sarah, relCode)
	if err == nil {
		t.Errorf(`NewRelatedName(Sarah, %q) succeeded, should have failed`,
			relCode)
	} else if e, ok := err.(*relname.BadRelatorCode); !ok {
		t.Errorf(`NewRelatedName(Sarah, %q) → wierd error %#v`,
			relCode, err)
	} else if *e != (relname.BadRelatorCode{Sarah, relCode}) {
		t.Errorf(`NewRelatedName(Sarah, %q) → error %#v\n\t%s %q)`,
			relCode, err,
			`right type but expected &BadRelatorCode{"Sarah A. Hoyt", `,
			relCode)
	}
}

func expectRelName(t *testing.T, method, actual, expected string) {
	if actual != expected {
		t.Errorf("expected %-16q but got %-16q from SarahA.%s()\n",
			expected, actual, method)
	}
}
