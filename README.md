# relname
Package **relname** handles names of people and organizations. It was written for working with metadata for ebooks.

Provides both 'common' (eg., "Dave Freer") and 'file-as' ("Freer, Dave") versions of names. Has special provisions for one-phrase names (eg., "Baen Books") and people with [generational titles](https://en.wikipedia.org/wiki/Suffix_(name)#Generational_titles). Assumes English-language naming conventions.

Provides two types, Name and RelName. The latter is a Name plus a [MARC relator code](https://www.loc.gov/marc/relators/relaterm.html) to specify the relationship between a named person or organization and the ebook (or any other creative work).
