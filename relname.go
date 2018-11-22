// File:	ebooks/relname/relname.go
// By:		Chris Chittleborough (github.com/c12h), November 2018

// Package “relname” provides types denoting the name of a person or
// organization in relation to ebooks.  It is intended for use by people who
// speak (or at least read) English.
//
//
// Kinds of Name
//
// This package supports both the ‘common’ form of names (eg., “Larry Correia”)
// and the ‘file-as’ form (“Correia, Larry”).  To cope with differences in
// converting to and from file-as format, “name” supports three kinds of names:
// one-part names, two-part names and (you guessed it) three-part names.
//
// Organizations have one part-names, possibly containing multiple words (eg.,
// “Baen Books”).  Some people also use one-part names (eg., “Teller”, “Gharlane
// of Eddore”; both are in Wikipedia).  For one-part names, the common and
// file-as forms are identical.
//
// Most people in the Anglosphere have two-part names: given names plus a
// surname (AKA family name).  Generally, people use one given name, omitting or
// abbreviating the others (eg., “Robert A. Heinlein”, “Mark L. Van Name”); some
// abbreviate them all (eg., “P. C. Hodgell”).  The distinctive thing about
// these names is that the file-as format is the last word or words from the
// common format, followed by a comma and the rest of the common name:
//	Robert A. Heinlein	←→	Heinlein, Robert A.
//	Mark L. Van Name	←→	Van Name, Mark L.
//	P. C. Hodgell		←→	Hodgell, P. C.
// For lack of a better term, we (ab)use the term ‘forename’ in this package to
// mean the part of a name before the surname.
//
// Some names have a third part, the generation (a suffix such as “Jr”, “Sr”,
// “père”, “fils” or a roman numeral), as well as a forename and surname.  See
// “https://en.wikipedia.org/wiki/Suffix_(name)#Generational_titles” for more
// details about these ‘generational’ names.  For example:
//	James Tiptree Jr.	←→	Tiptree, James Jr.
// Unlike forenames and surnames, generation parts never contain spaces.
//
// All strings passed to NewName1(), NewName2() and NewName3() are normalized by
// CleanString(), a utility function, before use.  This strips leading and
// trailing whitespace and replaces any internal whitespace sequences by a
// single U+0020 character.
//
//
// Related Names
//
// This package also supports ‘related names’. A RelatedName is a Name with a
// 3-letter MARC ‘relator’ code specifying the connection between a creative
// work and a person or organization.  (MARC is the MAchine-Readable Cataloging
// standard from the U.S. Library of Congress.)  All Relator codes consist of
// three letters from a-z: “aut” for author, “edc” for editor of compilation,
// etc.  See “https://www.loc.gov/marc/relators/relaterm.html” for the (very!)
// complete list.
//
//
// Methods
//
// Here’s a table showing the methods on Name and RelatedName and what they
// return for the various kinds of name, including the zero value (in the first column).
//
//	Common()	""	"Baen Books"	"Dave Freer"	"James Tiptree Jr."
//	FileAs()	""	"Baen Books"	"Freer, Dave"	"Tiptree, James Jr."
//	Surname()	""	"Baen Books"	"Freer"		"Tiptree"
//	Forename()	""	""		"Dave"		"James"
//	Generation()	""	""		""		"Jr."
//	NumParts()	0	1		2		3
//	String()	""	"Baen Books"	"Dave Freer"	"James Tiptree Jr."
//
//
// Limitations
//
// This package does not support post-nomials (“Ph.D”, “Esquire”, “FRS”, etc).
// Neither does it really support title prefixes (“Mr”, “Dr”, “Professor”, “Sir”
// etc), which it will treat as part of the forename; in some cases, that may be
// good enough.
//
//
package relname

// Wikipedia’articles on “personal name”, “surname”, “given name”, and “suffix
// (name)” may be useful, not least for showing the wide variety of naming
// conventions that this package does not support.
//
// That last article’s section on “Generational titles” is especially relevant.
// It mentions the following suffixes:
//	USA:     	Sr.		Jr.  II		III		IV ...
//	Britain: 	Snr		Jnr
//	France:		père		fils
//
// A regular expression resembling
//	\  ( Sr\. | \Jr. | (?: (?: L?X{0,3})? (?: I[XV] | V?I{0,3} ) ) ) $
// (with all but the first space removed) might be useful here.

//
// By the way, the authors whose names are used in comments in this file are
// all highly recommended.
//

import (
	"fmt"
	"regexp"
)

var reWhitespace = regexp.MustCompile(`\s+`)

// CleanString is a utility function which returns a string with any sequence of
// one or more whitespace characters (' ', '\t', '\n', '\r' etc etc) replaced by
// one space (' '), and all leading or trailing whitespace removed.
func CleanString(s string) string {
	//D// fmt.Printf(" #D# cleaning %q", s) //D//
	s = reWhitespace.ReplaceAllString(s, " ")
	//D// fmt.Printf(" → %q", s) //D//
	first, last := 0, len(s)

	if last > 0 {
		if s[0] == ' ' {
			if last == 1 {
				return ""
			}
			first++
		}
		if s[last-1] == ' ' {
			last--
		}
	}
	return s[first:last]
}

/*=============================== Name objects ===============================*/

// A Name holds the the name of a person or organization.
type Name struct {
	text      string
	boSurname uint16 // index of first byte of surname
	eoSurname uint16 // 1 + index of last byte of surname
}

// Invariants:	text !~ /^\s/, !~ /\s$/, !~ /\s\s/
//		0 <= boSurname <= eoSurname <= len(text)
//		boSurname == 0		 ⇒  eoSurname == len(text)
//		eoSurname != len(text) ⇒  eoSurname < len(text) - 1
//
//	For ‘simple’ names:		boSurname == 0, eoSurname == len(text)
//	For ‘typical’ names:		boSurname > 0,  eoSurname == len(text)
//	For ‘generational’ names:	boSurname > 0,  eoSurname < len(text)

// NewName1 constructs a one-part name.  Use it for organizations.
// The argument must contain at least one non-whitespace character, or NewName1 will
// return an error (and a zero-valued Name object).
func NewName1(text string) (Name, error) {
	t := CleanString(text)
	if t == "" {
		return Name{}, &EmptyPartError{1, text, "", ""}
	}
	return Name{text, 0, uint16(len(text))}, nil
}

// NewName2 constructs a two-part name.  Use it for most people.
// It returns an error (and an zero-valued Name object) if either argument is
// empty or contains only whitespace.
func NewName2(forename, surname string) (Name, error) {
	f := CleanString(forename)
	s := CleanString(surname)
	if f == "" || s == "" {
		return Name{}, &EmptyPartError{2, forename, surname, ""}
		//	f != "", s != "", false, forename, surname)
	}
	text := f + " " + s
	return Name{text, uint16(len(f) + 1), uint16(len(text))}, nil
}

// NewName3 constructs a three-part name.  Use it for people with generational
// suffixes.  It returns an error (and an zero-valued Name object) if any
// argument is empty or whitespace-only.
func NewName3(forename, surname, generation string) (Name, error) {
	f := CleanString(forename)
	s := CleanString(surname)
	g := CleanString(generation)
	if f == "" || s == "" || g == "" {
		return Name{}, &EmptyPartError{3, forename, surname, generation}
	}
	text := f + " " + s
	return Name{text + " " + g, uint16(len(f) + 1), uint16(len(text))}, nil
}

// Common returns the common (as opposed to file-as) form of a name.  For a
// zero-valued Name, it returns "".
func (n Name) Common() string { return n.text }

// FileAs returns the ‘file-as form’ of a name (eg., "Drake, David" rather than
// "David Drake").
func (n Name) FileAs() string {
	if n.boSurname == 0 {
		return n.text
	}
	faName := n.text[n.boSurname:n.eoSurname] + ", " + n.text[:n.boSurname-1]
	if n.eoSurname < uint16(len(n.text)) {
		faName += " " + n.text[n.eoSurname+1:]
	}
	return faName
}

// Surname returns the main part of a name.  It returns an empty string if and
// only if called on a zero-valued Name object.  Remember that surnames can
// contain multiple words.
func (n Name) Surname() string { return n.text[n.boSurname:n.eoSurname] }

// Forename returns the part of a person’s name that usually comes before the surname.
// It returns an empty string for zero and one-part names.
func (n Name) Forename() string {
	if n.boSurname == 0 {
		return ""
	}
	return n.text[:n.boSurname-1]
}

// Generation returns the generational suffix of a person’s name, or "" if this
// is not a three-part name.
func (n Name) Generation() string {
	if n.eoSurname == uint16(len(n.text)) {
		return ""
	}
	return n.text[n.eoSurname+1:]
}

// NumParts reports whether a name is one-part, two-part or three-part by returning
// 1, 2 or 3.  It returns 0 for zero-valued names.
func (n Name) NumParts() int {
	if n.boSurname == 0 {
		if n.text == "" {
			return 0
		} else {
			return 1
		}
	} else if n.eoSurname == uint16(len(n.text)) {
		return 2
	} else {
		return 3
	}
}

// String implements the fmt.Stringer interface. It returns the same value as Common().
func (n Name) String() string { return n.text }

/*=========================== RelatedName objects ============================*/

// A RelatedName is a Name plus a 3-letter ‘relator’ code; all three letters will
// be from a-z (no accents, never æ, þ etc).  (The U.S. Library of Congress has
// a list of relator codes at https://www.loc.gov/marc/relators/relaterm.html.)
type RelatedName struct {
	Name
	relCode [3]byte
}

var reRelator = regexp.MustCompile(`^[a-z][a-z][a-z]$`)

// NewRelatedName forms a RelatedName object by copying a Name object (which must not
// be zero-valued) and a relator code (which must consist of 3 letters in a-z).
func NewRelatedName(n Name, relatorCode string) (RelatedName, error) {
	if !reRelator.MatchString(relatorCode) {
		return RelatedName{}, &BadRelatorCode{n, relatorCode}
	}
	if n.NumParts() == 0 {
		return RelatedName{}, &BadName{relatorCode}
	}
	var code [3]byte
	code[0], code[1], code[2] = relatorCode[0], relatorCode[1], relatorCode[2]
	return RelatedName{Name: n, relCode: code}, nil
}

// Relator returns the three-letter relator code from a related name.
func (rn RelatedName) Relator() string {
	b := []byte("!!!")
	b[0], b[1], b[2] = rn.relCode[0], rn.relCode[1], rn.relCode[2]
	return string(b)
}

// String implements the fmt.Stringer interface.
func (rn RelatedName) String() string { return rn.Name.text + " (" + rn.Relator() + ")" }

/*================================== Errors ==================================*/

// EmptyPartError reports that NewName1(), NewName2() or NewName3() was given a
// name part that was empty or contained only whitespace characters.
type EmptyPartError struct {
	NumArgs          int
	Arg1, Arg2, Arg3 string
}

func (epe *EmptyPartError) Error() string {
	nEmpty, arg1, arg2, arg3 := 0, "", "", ""
	switch epe.NumArgs {
	default:
		if CleanString(epe.Arg3) == "" {
			nEmpty++
		}
		arg3 = fmt.Sprintf(", %q", epe.Arg3)
		fallthrough
	case 2:
		if CleanString(epe.Arg2) == "" {
			nEmpty++
		}
		arg2 = fmt.Sprintf(", %q", epe.Arg2)
		fallthrough
	case 1:
		if CleanString(epe.Arg1) == "" {
			nEmpty++
		}
		arg1 = fmt.Sprintf("%q", epe.Arg1)
	case 0:
		return fmt.Sprintf("BUG: bad EmptyPartError value %#v", *epe)
	}
	plural := "s"
	if nEmpty == 1 {
		plural = ""
	}
	return fmt.Sprintf("empty or whitespace-only argument%s in NewName%d(%s%s%s)",
		plural, epe.NumArgs, arg1, arg2, arg3)
}

// BadRelatorCode reports that NewRelatedName() was given an invalid relator code.
type BadRelatorCode struct {
	N Name
	C string
}

func (brc *BadRelatorCode) Error() string {
	return fmt.Sprintf("NewRelatedName(%q,%q): need /^[a-z][a-z][a-z]$/ for 2nd arg",
		brc.N, brc.C)
}

// BadName reports that NewRelatedName() was given a zero-valued Name argument.
type BadName struct {
	C string
}

func (bn *BadName) Error() string {
	return fmt.Sprintf("NewRelatedName(Name{},%q): need a non-zero-value Name", bn.C)
}
