7. The Language Specification - Variables
*****************************************

The next few chapters will be formal specifications about the language. Each chapter will have listings in alphabetical order. Since many function have completely unrelated uses, each will be listed by its primary use. This chapter is about the preinitialised variables in Pyth.

7.1. "G" - The Alphabet
=======================

This variable is preinitialised to the lowercase letters in the alphabet (i.e. "abcdefghijklmnopqrstuvwxyz"). 

Ex::

	==================================================
	@G5
	==================================================
	Pprint("\n",lookup(G,5))
	==================================================
	f

7.2. "H" - Empty Dictionary
===========================

This variable is set to an empty Python dictionary. Pyth also has a dictionary constructor, ``.d``.

Ex::

	==================================================
	XH"a"5
	==================================================
	Pprint("\n",assign_at(H,"a",5))
	==================================================
	{'a': 5}

7.3. "J" - Auto-Assignment With Copy
====================================

This, like ``K`` gets auto assigned the first time it is used. However, this is not directly assigned but assigned to a copy.

Ex::

	==================================================
	J5*3J
	==================================================
	J=copy(5)
	Pprint("\n",times(3,J))
	==================================================
	15

7.4. "K" - Auto-Assignment
==========================

The first time time this variable is mentioned, it assigns itself to the next expression. Unlike ``J``, this is not assigned to a copy but instead directly. The difference is relevant for mutable data types.

Ex::

	==================================================
	K7+TK
	==================================================
	K=7
	Pprint("\n",plus(T,K))
	==================================================
	17

7.5. "N" - Double Quote
=======================

This is pre-set to a string containing only a double quote. This useful since its one character shorter than ``\"``.

Ex::

	==================================================
	+++"Jane said "N"Hello!"N
	==================================================
	Pprint("\n",plus(plus(plus("Jane said ",N),"Hello!"),N))
	==================================================
	Jane said "Hello!"

7.6. "Q" - Evaluated Input
==========================

This variable auto-initializes to the evaluated input. The parser checks whether ``Q`` is in the code, and if it is, adds a line to the top setting ``Q`` equal to the evaluated input. This is the primary form of input in most programs.

Ex::

	input: 10
	
	==================================================
	yQ
	==================================================
	Q=copy(literal_eval(input()))
	Pprint("\n",subsets(Q))
	==================================================
	20

7.7. "T" - Ten
==============

Pretty self-explanatory. It starts off equalling ten. Ten is a very useful value.

Ex::

	==================================================
	^T6
	==================================================
	Pprint("\n",Ppow(T,6))
	==================================================
	1000000

7.8. "Y" - Empty List
=====================

Just an empty list that comes in handy when appending throughout a loop.

Ex::

	==================================================
	lY
	==================================================
	Pprint("\n",Plen(Y))
	==================================================
	0

7.9. "Z" - Zero
===============

This starts of as another very useful value, 0.

Ex::

	==================================================
	*Z5
	==================================================
	Pprint("\n",times(Z,5))
	==================================================
	0

7.10. "b" - Line Break
======================

This is set to a newline character.

Ex::

	==================================================
	jbUT
	==================================================
	Pprint("\n",join(b,urange(T)))
	==================================================
	0
	1
	2
	3
	4
	5
	6
	7
	8
	9

7.11. "d" - Space
=================

This is set to a string containing a single space.

Ex::

	==================================================
	jdUT
	==================================================
	Pprint("\n",join(d,urange(T)))
	==================================================
	0 1 2 3 4 5 6 7 8 9

7.12. "k" - Empty String
========================

Pre-initialised to an empty string. Useful for joining.

Ex::

	==================================================
	jkUT
	==================================================
	Pprint("\n",join(k,urange(T)))
	==================================================
	0123456789

7.13. "z" - Raw Input
=====================

This is set to the input, like ``Q``, but not evaluated. This is useful for string input.

Ex::

	input: Hello
	
	==================================================
	*z5
	==================================================
	z=copy(input())
	Pprint("\n",times(z,5))
	==================================================
	HelloHelloHelloHelloHello
