package query

import (
	"fmt"
	"math"
	"sort"
	"strconv"
)

const endSymbol rune = 1114112

/* The rule types inferred from the grammar are below. */
type pegRule uint8

const (
	ruleUnknown pegRule = iota
	ruleQuery
	ruleJoin
	ruleJoinKey
	ruleJoinRow
	ruleKeyJoin
	ruleValueJoin
	ruleSelect
	ruleWherePart
	ruleSelectKey
	ruleLimit
	ruleCryptoKey
	ruleWhere
	ruleWhereClause
	ruleAndClause
	ruleOrClause
	rulePredicateClause
	rulePredicate
	rulePredicateValue
	rulePredicateRowKey
	rulePredicateKey
	rulePredicateLiteralValue
	ruleLiteral
	rulePositiveInteger
	ruleKey
	ruleKeySymbols
	ruleEscape
	ruleMustSpacing
	ruleSpacing
	ruleAction0
	ruleAction1
	rulePegText
	ruleAction2
	ruleAction3
	ruleAction4
	ruleAction5
	ruleAction6
	ruleAction7
	ruleAction8
	ruleAction9
	ruleAction10
	ruleAction11
	ruleAction12
	ruleAction13
	ruleAction14
	ruleAction15
	ruleAction16
	ruleAction17
	ruleAction18
)

var rul3s = [...]string{
	"Unknown",
	"Query",
	"Join",
	"JoinKey",
	"JoinRow",
	"KeyJoin",
	"ValueJoin",
	"Select",
	"WherePart",
	"SelectKey",
	"Limit",
	"CryptoKey",
	"Where",
	"WhereClause",
	"AndClause",
	"OrClause",
	"PredicateClause",
	"Predicate",
	"PredicateValue",
	"PredicateRowKey",
	"PredicateKey",
	"PredicateLiteralValue",
	"Literal",
	"PositiveInteger",
	"Key",
	"KeySymbols",
	"Escape",
	"MustSpacing",
	"Spacing",
	"Action0",
	"Action1",
	"PegText",
	"Action2",
	"Action3",
	"Action4",
	"Action5",
	"Action6",
	"Action7",
	"Action8",
	"Action9",
	"Action10",
	"Action11",
	"Action12",
	"Action13",
	"Action14",
	"Action15",
	"Action16",
	"Action17",
	"Action18",
}

type token32 struct {
	pegRule
	begin, end uint32
}

func (t *token32) String() string {
	return fmt.Sprintf("\x1B[34m%v\x1B[m %v %v", rul3s[t.pegRule], t.begin, t.end)
}

type node32 struct {
	token32
	up, next *node32
}

func (node *node32) print(pretty bool, buffer string) {
	var print func(node *node32, depth int)
	print = func(node *node32, depth int) {
		for node != nil {
			for c := 0; c < depth; c++ {
				fmt.Printf(" ")
			}
			rule := rul3s[node.pegRule]
			quote := strconv.Quote(string(([]rune(buffer)[node.begin:node.end])))
			if !pretty {
				fmt.Printf("%v %v\n", rule, quote)
			} else {
				fmt.Printf("\x1B[34m%v\x1B[m %v\n", rule, quote)
			}
			if node.up != nil {
				print(node.up, depth+1)
			}
			node = node.next
		}
	}
	print(node, 0)
}

func (node *node32) Print(buffer string) {
	node.print(false, buffer)
}

func (node *node32) PrettyPrint(buffer string) {
	node.print(true, buffer)
}

type tokens32 struct {
	tree []token32
}

func (t *tokens32) Trim(length uint32) {
	t.tree = t.tree[:length]
}

func (t *tokens32) Print() {
	for _, token := range t.tree {
		fmt.Println(token.String())
	}
}

func (t *tokens32) AST() *node32 {
	type element struct {
		node *node32
		down *element
	}
	tokens := t.Tokens()
	var stack *element
	for _, token := range tokens {
		if token.begin == token.end {
			continue
		}
		node := &node32{token32: token}
		for stack != nil && stack.node.begin >= token.begin && stack.node.end <= token.end {
			stack.node.next = node.up
			node.up = stack.node
			stack = stack.down
		}
		stack = &element{node: node, down: stack}
	}
	if stack != nil {
		return stack.node
	}
	return nil
}

func (t *tokens32) PrintSyntaxTree(buffer string) {
	t.AST().Print(buffer)
}

func (t *tokens32) PrettyPrintSyntaxTree(buffer string) {
	t.AST().PrettyPrint(buffer)
}

func (t *tokens32) Add(rule pegRule, begin, end, index uint32) {
	if tree := t.tree; int(index) >= len(tree) {
		expanded := make([]token32, 2*len(tree))
		copy(expanded, tree)
		t.tree = expanded
	}
	t.tree[index] = token32{
		pegRule: rule,
		begin:   begin,
		end:     end,
	}
}

func (t *tokens32) Tokens() []token32 {
	return t.tree
}

type QueryParser struct {
	QueryAST

	Buffer string
	buffer []rune
	rules  [49]func() bool
	parse  func(rule ...int) error
	reset  func()
	Pretty bool
	tokens32
}

func (p *QueryParser) Parse(rule ...int) error {
	return p.parse(rule...)
}

func (p *QueryParser) Reset() {
	p.reset()
}

type textPosition struct {
	line, symbol int
}

type textPositionMap map[int]textPosition

func translatePositions(buffer []rune, positions []int) textPositionMap {
	length, translations, j, line, symbol := len(positions), make(textPositionMap, len(positions)), 0, 1, 0
	sort.Ints(positions)

search:
	for i, c := range buffer {
		if c == '\n' {
			line, symbol = line+1, 0
		} else {
			symbol++
		}
		if i == positions[j] {
			translations[positions[j]] = textPosition{line, symbol}
			for j++; j < length; j++ {
				if i != positions[j] {
					continue search
				}
			}
			break search
		}
	}

	return translations
}

type parseError struct {
	p   *QueryParser
	max token32
}

func (e *parseError) Error() string {
	tokens, error := []token32{e.max}, "\n"
	positions, p := make([]int, 2*len(tokens)), 0
	for _, token := range tokens {
		positions[p], p = int(token.begin), p+1
		positions[p], p = int(token.end), p+1
	}
	translations := translatePositions(e.p.buffer, positions)
	format := "parse error near %v (line %v symbol %v - line %v symbol %v):\n%v\n"
	if e.p.Pretty {
		format = "parse error near \x1B[34m%v\x1B[m (line %v symbol %v - line %v symbol %v):\n%v\n"
	}
	for _, token := range tokens {
		begin, end := int(token.begin), int(token.end)
		error += fmt.Sprintf(format,
			rul3s[token.pegRule],
			translations[begin].line, translations[begin].symbol,
			translations[end].line, translations[end].symbol,
			strconv.Quote(string(e.p.buffer[begin:end])))
	}

	return error
}

func (p *QueryParser) PrintSyntaxTree() {
	if p.Pretty {
		p.tokens32.PrettyPrintSyntaxTree(p.Buffer)
	} else {
		p.tokens32.PrintSyntaxTree(p.Buffer)
	}
}

func (p *QueryParser) Execute() {
	buffer, _buffer, text, begin, end := p.Buffer, p.buffer, "", 0, 0
	for _, token := range p.Tokens() {
		switch token.pegRule {

		case rulePegText:
			begin, end = int(token.begin), int(token.end)
			text = string(_buffer[begin:end])

		case ruleAction0:
			p.AddSelect()
		case ruleAction1:
			p.AddJoin()
		case ruleAction2:
			p.SetTableName(buffer[begin:end])
		case ruleAction3:
			p.AddJoinRow()
		case ruleAction4:
			p.SetJoinRowKey(buffer[begin:end])
		case ruleAction5:
			p.SetJoinKey(buffer[begin:end])
		case ruleAction6:
			p.SetJoinValue(buffer[begin:end])
		case ruleAction7:
			p.SetTableName(buffer[begin:end])
		case ruleAction8:
			p.SetLimit(buffer[begin:end])
		case ruleAction9:
			p.AddCryptoKey(buffer[begin:end])
		case ruleAction10:
			p.PushWhere()
		case ruleAction11:
			p.PopWhere()
		case ruleAction12:
			p.SetWhereCommand("and")
		case ruleAction13:
			p.SetWhereCommand("or")
		case ruleAction14:
			p.InitPredicate()
		case ruleAction15:
			p.SetPredicateCommand(buffer[begin:end])
		case ruleAction16:
			p.UsePredicateRowKey()
		case ruleAction17:
			p.AddPredicateKey(buffer[begin:end])
		case ruleAction18:
			p.AddPredicateLiteral(buffer[begin:end])

		}
	}
	_, _, _, _, _ = buffer, _buffer, text, begin, end
}

func (p *QueryParser) Init() {
	var (
		max                  token32
		position, tokenIndex uint32
		buffer               []rune
	)
	p.reset = func() {
		max = token32{}
		position, tokenIndex = 0, 0

		p.buffer = []rune(p.Buffer)
		if len(p.buffer) == 0 || p.buffer[len(p.buffer)-1] != endSymbol {
			p.buffer = append(p.buffer, endSymbol)
		}
		buffer = p.buffer
	}
	p.reset()

	_rules := p.rules
	tree := tokens32{tree: make([]token32, math.MaxInt16)}
	p.parse = func(rule ...int) error {
		r := 1
		if len(rule) > 0 {
			r = rule[0]
		}
		matches := p.rules[r]()
		p.tokens32 = tree
		if matches {
			p.Trim(tokenIndex)
			return nil
		}
		return &parseError{p, max}
	}

	add := func(rule pegRule, begin uint32) {
		tree.Add(rule, begin, position, tokenIndex)
		tokenIndex++
		if begin != position && position > max.end {
			max = token32{rule, begin, position}
		}
	}

	matchDot := func() bool {
		if buffer[position] != endSymbol {
			position++
			return true
		}
		return false
	}

	/*matchChar := func(c byte) bool {
		if buffer[position] == c {
			position++
			return true
		}
		return false
	}*/

	/*matchRange := func(lower byte, upper byte) bool {
		if c := buffer[position]; c >= lower && c <= upper {
			position++
			return true
		}
		return false
	}*/

	_rules = [...]func() bool{
		nil,
		/* 0 Query <- <(Spacing ((Select Action0) / (Join Action1)) Spacing !.)> */
		func() bool {
			position0, tokenIndex0 := position, tokenIndex
			{
				position1 := position
				if !_rules[ruleSpacing]() {
					goto l0
				}
				{
					position2, tokenIndex2 := position, tokenIndex
					{
						position4 := position
						if buffer[position] != rune('s') {
							goto l3
						}
						position++
						if buffer[position] != rune('e') {
							goto l3
						}
						position++
						if buffer[position] != rune('l') {
							goto l3
						}
						position++
						if buffer[position] != rune('e') {
							goto l3
						}
						position++
						if buffer[position] != rune('c') {
							goto l3
						}
						position++
						if buffer[position] != rune('t') {
							goto l3
						}
						position++
						if !_rules[ruleMustSpacing]() {
							goto l3
						}
						{
							position5 := position
							{
								position6 := position
								if !_rules[ruleKey]() {
									goto l3
								}
								add(rulePegText, position6)
							}
							{
								add(ruleAction7, position)
							}
							add(ruleSelectKey, position5)
						}
					l8:
						{
							position9, tokenIndex9 := position, tokenIndex
							if !_rules[ruleMustSpacing]() {
								goto l9
							}
							{
								position10 := position
								{
									switch buffer[position] {
									case 's':
										if !_rules[ruleCryptoKey]() {
											goto l9
										}
										break
									case 'l':
										{
											position12 := position
											if buffer[position] != rune('l') {
												goto l9
											}
											position++
											if buffer[position] != rune('i') {
												goto l9
											}
											position++
											if buffer[position] != rune('m') {
												goto l9
											}
											position++
											if buffer[position] != rune('i') {
												goto l9
											}
											position++
											if buffer[position] != rune('t') {
												goto l9
											}
											position++
											if !_rules[ruleMustSpacing]() {
												goto l9
											}
											{
												position13 := position
												{
													position14 := position
													if c := buffer[position]; c < rune('1') || c > rune('9') {
														goto l9
													}
													position++
												l15:
													{
														position16, tokenIndex16 := position, tokenIndex
														if c := buffer[position]; c < rune('0') || c > rune('9') {
															goto l16
														}
														position++
														goto l15
													l16:
														position, tokenIndex = position16, tokenIndex16
													}
													add(rulePositiveInteger, position14)
												}
												add(rulePegText, position13)
											}
											{
												add(ruleAction8, position)
											}
											add(ruleLimit, position12)
										}
										break
									default:
										{
											position18 := position
											if buffer[position] != rune('w') {
												goto l9
											}
											position++
											if buffer[position] != rune('h') {
												goto l9
											}
											position++
											if buffer[position] != rune('e') {
												goto l9
											}
											position++
											if buffer[position] != rune('r') {
												goto l9
											}
											position++
											if buffer[position] != rune('e') {
												goto l9
											}
											position++
											if !_rules[ruleMustSpacing]() {
												goto l9
											}
											if !_rules[ruleWhereClause]() {
												goto l9
											}
											add(ruleWhere, position18)
										}
										break
									}
								}

								add(ruleWherePart, position10)
							}
							goto l8
						l9:
							position, tokenIndex = position9, tokenIndex9
						}
						add(ruleSelect, position4)
					}
					{
						add(ruleAction0, position)
					}
					goto l2
				l3:
					position, tokenIndex = position2, tokenIndex2
					{
						position20 := position
						if buffer[position] != rune('j') {
							goto l0
						}
						position++
						if buffer[position] != rune('o') {
							goto l0
						}
						position++
						if buffer[position] != rune('i') {
							goto l0
						}
						position++
						if buffer[position] != rune('n') {
							goto l0
						}
						position++
						if !_rules[ruleMustSpacing]() {
							goto l0
						}
						{
							position21 := position
							{
								position22 := position
								if !_rules[ruleKey]() {
									goto l0
								}
								add(rulePegText, position22)
							}
							{
								add(ruleAction2, position)
							}
							add(ruleJoinKey, position21)
						}
					l24:
						{
							position25, tokenIndex25 := position, tokenIndex
							if !_rules[ruleMustSpacing]() {
								goto l25
							}
							if !_rules[ruleCryptoKey]() {
								goto l25
							}
							goto l24
						l25:
							position, tokenIndex = position25, tokenIndex25
						}
						if !_rules[ruleMustSpacing]() {
							goto l0
						}
						if buffer[position] != rune('r') {
							goto l0
						}
						position++
						if buffer[position] != rune('o') {
							goto l0
						}
						position++
						if buffer[position] != rune('w') {
							goto l0
						}
						position++
						if buffer[position] != rune('s') {
							goto l0
						}
						position++
						if !_rules[ruleMustSpacing]() {
							goto l0
						}
						if !_rules[ruleJoinRow]() {
							goto l0
						}
					l26:
						{
							position27, tokenIndex27 := position, tokenIndex
							if !_rules[ruleSpacing]() {
								goto l27
							}
							if buffer[position] != rune(',') {
								goto l27
							}
							position++
							if !_rules[ruleSpacing]() {
								goto l27
							}
							if !_rules[ruleJoinRow]() {
								goto l27
							}
							goto l26
						l27:
							position, tokenIndex = position27, tokenIndex27
						}
						if !_rules[ruleSpacing]() {
							goto l0
						}
						add(ruleJoin, position20)
					}
					{
						add(ruleAction1, position)
					}
				}
			l2:
				if !_rules[ruleSpacing]() {
					goto l0
				}
				{
					position29, tokenIndex29 := position, tokenIndex
					if !matchDot() {
						goto l29
					}
					goto l0
				l29:
					position, tokenIndex = position29, tokenIndex29
				}
				add(ruleQuery, position1)
			}
			return true
		l0:
			position, tokenIndex = position0, tokenIndex0
			return false
		},
		/* 1 Join <- <('j' 'o' 'i' 'n' MustSpacing JoinKey (MustSpacing CryptoKey)* MustSpacing ('r' 'o' 'w' 's') MustSpacing JoinRow (Spacing ',' Spacing JoinRow)* Spacing)> */
		nil,
		/* 2 JoinKey <- <(<Key> Action2)> */
		nil,
		/* 3 JoinRow <- <(Action3 '(' Spacing KeyJoin Spacing (',' Spacing ValueJoin Spacing)* ')')> */
		func() bool {
			position32, tokenIndex32 := position, tokenIndex
			{
				position33 := position
				{
					add(ruleAction3, position)
				}
				if buffer[position] != rune('(') {
					goto l32
				}
				position++
				if !_rules[ruleSpacing]() {
					goto l32
				}
				{
					position35 := position
					if buffer[position] != rune('@') {
						goto l32
					}
					position++
					if buffer[position] != rune('k') {
						goto l32
					}
					position++
					if buffer[position] != rune('e') {
						goto l32
					}
					position++
					if buffer[position] != rune('y') {
						goto l32
					}
					position++
					if !_rules[ruleSpacing]() {
						goto l32
					}
					if buffer[position] != rune('=') {
						goto l32
					}
					position++
					if !_rules[ruleSpacing]() {
						goto l32
					}
					{
						position36, tokenIndex36 := position, tokenIndex
						if buffer[position] != rune('@') {
							goto l37
						}
						position++
						if buffer[position] != rune('"') {
							goto l37
						}
						position++
						{
							position38 := position
							if !_rules[ruleLiteral]() {
								goto l37
							}
							add(rulePegText, position38)
						}
						if buffer[position] != rune('"') {
							goto l37
						}
						position++
						goto l36
					l37:
						position, tokenIndex = position36, tokenIndex36
						{
							position39 := position
							if !_rules[ruleKey]() {
								goto l32
							}
							add(rulePegText, position39)
						}
					}
				l36:
					{
						add(ruleAction4, position)
					}
					add(ruleKeyJoin, position35)
				}
				if !_rules[ruleSpacing]() {
					goto l32
				}
			l41:
				{
					position42, tokenIndex42 := position, tokenIndex
					if buffer[position] != rune(',') {
						goto l42
					}
					position++
					if !_rules[ruleSpacing]() {
						goto l42
					}
					{
						position43 := position
						{
							position44, tokenIndex44 := position, tokenIndex
							{
								position46 := position
								if !_rules[ruleKey]() {
									goto l45
								}
								add(rulePegText, position46)
							}
							goto l44
						l45:
							position, tokenIndex = position44, tokenIndex44
							if buffer[position] != rune('@') {
								goto l42
							}
							position++
							if buffer[position] != rune('"') {
								goto l42
							}
							position++
							{
								position47 := position
								if !_rules[ruleLiteral]() {
									goto l42
								}
								add(rulePegText, position47)
							}
							if buffer[position] != rune('"') {
								goto l42
							}
							position++
						}
					l44:
						{
							add(ruleAction5, position)
						}
						if !_rules[ruleSpacing]() {
							goto l42
						}
						if buffer[position] != rune('=') {
							goto l42
						}
						position++
						if !_rules[ruleSpacing]() {
							goto l42
						}
						if buffer[position] != rune('"') {
							goto l42
						}
						position++
						{
							position49 := position
							if !_rules[ruleLiteral]() {
								goto l42
							}
							add(rulePegText, position49)
						}
						if buffer[position] != rune('"') {
							goto l42
						}
						position++
						{
							add(ruleAction6, position)
						}
						add(ruleValueJoin, position43)
					}
					if !_rules[ruleSpacing]() {
						goto l42
					}
					goto l41
				l42:
					position, tokenIndex = position42, tokenIndex42
				}
				if buffer[position] != rune(')') {
					goto l32
				}
				position++
				add(ruleJoinRow, position33)
			}
			return true
		l32:
			position, tokenIndex = position32, tokenIndex32
			return false
		},
		/* 4 KeyJoin <- <('@' 'k' 'e' 'y' Spacing '=' Spacing (('@' '"' <Literal> '"') / <Key>) Action4)> */
		nil,
		/* 5 ValueJoin <- <((<Key> / ('@' '"' <Literal> '"')) Action5 Spacing '=' Spacing '"' <Literal> '"' Action6)> */
		nil,
		/* 6 Select <- <('s' 'e' 'l' 'e' 'c' 't' MustSpacing SelectKey (MustSpacing WherePart)*)> */
		nil,
		/* 7 WherePart <- <((&('s') CryptoKey) | (&('l') Limit) | (&('w') Where))> */
		nil,
		/* 8 SelectKey <- <(<Key> Action7)> */
		nil,
		/* 9 Limit <- <('l' 'i' 'm' 'i' 't' MustSpacing <PositiveInteger> Action8)> */
		nil,
		/* 10 CryptoKey <- <('s' 'i' 'g' 'n' 'e' 'd' MustSpacing '"' <KeySymbols> '"' Action9)> */
		func() bool {
			position57, tokenIndex57 := position, tokenIndex
			{
				position58 := position
				if buffer[position] != rune('s') {
					goto l57
				}
				position++
				if buffer[position] != rune('i') {
					goto l57
				}
				position++
				if buffer[position] != rune('g') {
					goto l57
				}
				position++
				if buffer[position] != rune('n') {
					goto l57
				}
				position++
				if buffer[position] != rune('e') {
					goto l57
				}
				position++
				if buffer[position] != rune('d') {
					goto l57
				}
				position++
				if !_rules[ruleMustSpacing]() {
					goto l57
				}
				if buffer[position] != rune('"') {
					goto l57
				}
				position++
				{
					position59 := position
					if !_rules[ruleKeySymbols]() {
						goto l57
					}
					add(rulePegText, position59)
				}
				if buffer[position] != rune('"') {
					goto l57
				}
				position++
				{
					add(ruleAction9, position)
				}
				add(ruleCryptoKey, position58)
			}
			return true
		l57:
			position, tokenIndex = position57, tokenIndex57
			return false
		},
		/* 11 Where <- <('w' 'h' 'e' 'r' 'e' MustSpacing WhereClause)> */
		nil,
		/* 12 WhereClause <- <(Action10 ((&('s') PredicateClause) | (&('o') OrClause) | (&('a') AndClause)) Action11)> */
		func() bool {
			position62, tokenIndex62 := position, tokenIndex
			{
				position63 := position
				{
					add(ruleAction10, position)
				}
				{
					switch buffer[position] {
					case 's':
						{
							position66 := position
							{
								add(ruleAction14, position)
							}
							{
								position68 := position
								{
									position69 := position
									if buffer[position] != rune('s') {
										goto l62
									}
									position++
									if buffer[position] != rune('t') {
										goto l62
									}
									position++
									if buffer[position] != rune('r') {
										goto l62
									}
									position++
									if buffer[position] != rune('_') {
										goto l62
									}
									position++
									if buffer[position] != rune('e') {
										goto l62
									}
									position++
									if buffer[position] != rune('q') {
										goto l62
									}
									position++
									add(rulePegText, position69)
								}
								{
									add(ruleAction15, position)
								}
								add(rulePredicate, position68)
							}
							if !_rules[ruleSpacing]() {
								goto l62
							}
							if buffer[position] != rune('(') {
								goto l62
							}
							position++
							if !_rules[ruleSpacing]() {
								goto l62
							}
							if !_rules[rulePredicateValue]() {
								goto l62
							}
						l71:
							{
								position72, tokenIndex72 := position, tokenIndex
								if buffer[position] != rune(',') {
									goto l72
								}
								position++
								if !_rules[ruleSpacing]() {
									goto l72
								}
								if !_rules[rulePredicateValue]() {
									goto l72
								}
								if !_rules[ruleSpacing]() {
									goto l72
								}
								goto l71
							l72:
								position, tokenIndex = position72, tokenIndex72
							}
							if buffer[position] != rune(')') {
								goto l62
							}
							position++
							add(rulePredicateClause, position66)
						}
						break
					case 'o':
						{
							position73 := position
							if buffer[position] != rune('o') {
								goto l62
							}
							position++
							if buffer[position] != rune('r') {
								goto l62
							}
							position++
							{
								add(ruleAction13, position)
							}
							if !_rules[ruleSpacing]() {
								goto l62
							}
							if buffer[position] != rune('(') {
								goto l62
							}
							position++
							if !_rules[ruleSpacing]() {
								goto l62
							}
							if !_rules[ruleWhereClause]() {
								goto l62
							}
							if !_rules[ruleSpacing]() {
								goto l62
							}
						l75:
							{
								position76, tokenIndex76 := position, tokenIndex
								if buffer[position] != rune(',') {
									goto l76
								}
								position++
								if !_rules[ruleSpacing]() {
									goto l76
								}
								if !_rules[ruleWhereClause]() {
									goto l76
								}
								if !_rules[ruleSpacing]() {
									goto l76
								}
								goto l75
							l76:
								position, tokenIndex = position76, tokenIndex76
							}
							if buffer[position] != rune(')') {
								goto l62
							}
							position++
							add(ruleOrClause, position73)
						}
						break
					default:
						{
							position77 := position
							if buffer[position] != rune('a') {
								goto l62
							}
							position++
							if buffer[position] != rune('n') {
								goto l62
							}
							position++
							if buffer[position] != rune('d') {
								goto l62
							}
							position++
							{
								add(ruleAction12, position)
							}
							if !_rules[ruleSpacing]() {
								goto l62
							}
							if buffer[position] != rune('(') {
								goto l62
							}
							position++
							if !_rules[ruleSpacing]() {
								goto l62
							}
							if !_rules[ruleWhereClause]() {
								goto l62
							}
							if !_rules[ruleSpacing]() {
								goto l62
							}
						l79:
							{
								position80, tokenIndex80 := position, tokenIndex
								if buffer[position] != rune(',') {
									goto l80
								}
								position++
								if !_rules[ruleSpacing]() {
									goto l80
								}
								if !_rules[ruleWhereClause]() {
									goto l80
								}
								if !_rules[ruleSpacing]() {
									goto l80
								}
								goto l79
							l80:
								position, tokenIndex = position80, tokenIndex80
							}
							if buffer[position] != rune(')') {
								goto l62
							}
							position++
							add(ruleAndClause, position77)
						}
						break
					}
				}

				{
					add(ruleAction11, position)
				}
				add(ruleWhereClause, position63)
			}
			return true
		l62:
			position, tokenIndex = position62, tokenIndex62
			return false
		},
		/* 13 AndClause <- <('a' 'n' 'd' Action12 Spacing '(' Spacing WhereClause Spacing (',' Spacing WhereClause Spacing)* ')')> */
		nil,
		/* 14 OrClause <- <('o' 'r' Action13 Spacing '(' Spacing WhereClause Spacing (',' Spacing WhereClause Spacing)* ')')> */
		nil,
		/* 15 PredicateClause <- <(Action14 Predicate Spacing '(' Spacing PredicateValue (',' Spacing PredicateValue Spacing)* ')')> */
		nil,
		/* 16 Predicate <- <(<('s' 't' 'r' '_' 'e' 'q')> Action15)> */
		nil,
		/* 17 PredicateValue <- <(PredicateRowKey / PredicateKey / PredicateLiteralValue)> */
		func() bool {
			position86, tokenIndex86 := position, tokenIndex
			{
				position87 := position
				{
					position88, tokenIndex88 := position, tokenIndex
					{
						position90 := position
						if buffer[position] != rune('@') {
							goto l89
						}
						position++
						if buffer[position] != rune('k') {
							goto l89
						}
						position++
						if buffer[position] != rune('e') {
							goto l89
						}
						position++
						if buffer[position] != rune('y') {
							goto l89
						}
						position++
						{
							add(ruleAction16, position)
						}
						add(rulePredicateRowKey, position90)
					}
					goto l88
				l89:
					position, tokenIndex = position88, tokenIndex88
					{
						position93 := position
						{
							position94, tokenIndex94 := position, tokenIndex
							{
								position96 := position
								if !_rules[ruleKey]() {
									goto l95
								}
								add(rulePegText, position96)
							}
							goto l94
						l95:
							position, tokenIndex = position94, tokenIndex94
							if buffer[position] != rune('@') {
								goto l92
							}
							position++
							if buffer[position] != rune('"') {
								goto l92
							}
							position++
							{
								position97 := position
								if !_rules[ruleLiteral]() {
									goto l92
								}
								add(rulePegText, position97)
							}
							if buffer[position] != rune('"') {
								goto l92
							}
							position++
						}
					l94:
						{
							add(ruleAction17, position)
						}
						add(rulePredicateKey, position93)
					}
					goto l88
				l92:
					position, tokenIndex = position88, tokenIndex88
					{
						position99 := position
						if buffer[position] != rune('"') {
							goto l86
						}
						position++
						{
							position100 := position
							if !_rules[ruleLiteral]() {
								goto l86
							}
							add(rulePegText, position100)
						}
						if buffer[position] != rune('"') {
							goto l86
						}
						position++
						{
							add(ruleAction18, position)
						}
						add(rulePredicateLiteralValue, position99)
					}
				}
			l88:
				add(rulePredicateValue, position87)
			}
			return true
		l86:
			position, tokenIndex = position86, tokenIndex86
			return false
		},
		/* 18 PredicateRowKey <- <('@' 'k' 'e' 'y' Action16)> */
		nil,
		/* 19 PredicateKey <- <((<Key> / ('@' '"' <Literal> '"')) Action17)> */
		nil,
		/* 20 PredicateLiteralValue <- <('"' <Literal> '"' Action18)> */
		nil,
		/* 21 Literal <- <(Escape / (!'"' .))*> */
		func() bool {
			{
				position106 := position
			l107:
				{
					position108, tokenIndex108 := position, tokenIndex
					{
						position109, tokenIndex109 := position, tokenIndex
						{
							position111 := position
							if buffer[position] != rune('\\') {
								goto l110
							}
							position++
							{
								switch buffer[position] {
								case 'v':
									if buffer[position] != rune('v') {
										goto l110
									}
									position++
									break
								case 't':
									if buffer[position] != rune('t') {
										goto l110
									}
									position++
									break
								case 'r':
									if buffer[position] != rune('r') {
										goto l110
									}
									position++
									break
								case 'n':
									if buffer[position] != rune('n') {
										goto l110
									}
									position++
									break
								case 'f':
									if buffer[position] != rune('f') {
										goto l110
									}
									position++
									break
								case 'b':
									if buffer[position] != rune('b') {
										goto l110
									}
									position++
									break
								case 'a':
									if buffer[position] != rune('a') {
										goto l110
									}
									position++
									break
								case '\\':
									if buffer[position] != rune('\\') {
										goto l110
									}
									position++
									break
								default:
									if buffer[position] != rune('"') {
										goto l110
									}
									position++
									break
								}
							}

							add(ruleEscape, position111)
						}
						goto l109
					l110:
						position, tokenIndex = position109, tokenIndex109
						{
							position113, tokenIndex113 := position, tokenIndex
							if buffer[position] != rune('"') {
								goto l113
							}
							position++
							goto l108
						l113:
							position, tokenIndex = position113, tokenIndex113
						}
						if !matchDot() {
							goto l108
						}
					}
				l109:
					goto l107
				l108:
					position, tokenIndex = position108, tokenIndex108
				}
				add(ruleLiteral, position106)
			}
			return true
		},
		/* 22 PositiveInteger <- <([1-9] [0-9]*)> */
		nil,
		/* 23 Key <- <KeySymbols> */
		func() bool {
			position115, tokenIndex115 := position, tokenIndex
			{
				position116 := position
				if !_rules[ruleKeySymbols]() {
					goto l115
				}
				add(ruleKey, position116)
			}
			return true
		l115:
			position, tokenIndex = position115, tokenIndex115
			return false
		},
		/* 24 KeySymbols <- <((&('-') '-') | (&('+') '+') | (&('.') '.') | (&('_') '_') | (&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') [0-9]) | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z') [A-Z]) | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [a-z]))+> */
		func() bool {
			position117, tokenIndex117 := position, tokenIndex
			{
				position118 := position
				{
					switch buffer[position] {
					case '-':
						if buffer[position] != rune('-') {
							goto l117
						}
						position++
						break
					case '+':
						if buffer[position] != rune('+') {
							goto l117
						}
						position++
						break
					case '.':
						if buffer[position] != rune('.') {
							goto l117
						}
						position++
						break
					case '_':
						if buffer[position] != rune('_') {
							goto l117
						}
						position++
						break
					case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l117
						}
						position++
						break
					case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
						if c := buffer[position]; c < rune('A') || c > rune('Z') {
							goto l117
						}
						position++
						break
					default:
						if c := buffer[position]; c < rune('a') || c > rune('z') {
							goto l117
						}
						position++
						break
					}
				}

			l119:
				{
					position120, tokenIndex120 := position, tokenIndex
					{
						switch buffer[position] {
						case '-':
							if buffer[position] != rune('-') {
								goto l120
							}
							position++
							break
						case '+':
							if buffer[position] != rune('+') {
								goto l120
							}
							position++
							break
						case '.':
							if buffer[position] != rune('.') {
								goto l120
							}
							position++
							break
						case '_':
							if buffer[position] != rune('_') {
								goto l120
							}
							position++
							break
						case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l120
							}
							position++
							break
						case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
							if c := buffer[position]; c < rune('A') || c > rune('Z') {
								goto l120
							}
							position++
							break
						default:
							if c := buffer[position]; c < rune('a') || c > rune('z') {
								goto l120
							}
							position++
							break
						}
					}

					goto l119
				l120:
					position, tokenIndex = position120, tokenIndex120
				}
				add(ruleKeySymbols, position118)
			}
			return true
		l117:
			position, tokenIndex = position117, tokenIndex117
			return false
		},
		/* 25 Escape <- <('\\' ((&('v') 'v') | (&('t') 't') | (&('r') 'r') | (&('n') 'n') | (&('f') 'f') | (&('b') 'b') | (&('a') 'a') | (&('\\') '\\') | (&('"') '"')))> */
		nil,
		/* 26 MustSpacing <- <((&('\n') '\n') | (&('\t') '\t') | (&(' ') ' '))+> */
		func() bool {
			position124, tokenIndex124 := position, tokenIndex
			{
				position125 := position
				{
					switch buffer[position] {
					case '\n':
						if buffer[position] != rune('\n') {
							goto l124
						}
						position++
						break
					case '\t':
						if buffer[position] != rune('\t') {
							goto l124
						}
						position++
						break
					default:
						if buffer[position] != rune(' ') {
							goto l124
						}
						position++
						break
					}
				}

			l126:
				{
					position127, tokenIndex127 := position, tokenIndex
					{
						switch buffer[position] {
						case '\n':
							if buffer[position] != rune('\n') {
								goto l127
							}
							position++
							break
						case '\t':
							if buffer[position] != rune('\t') {
								goto l127
							}
							position++
							break
						default:
							if buffer[position] != rune(' ') {
								goto l127
							}
							position++
							break
						}
					}

					goto l126
				l127:
					position, tokenIndex = position127, tokenIndex127
				}
				add(ruleMustSpacing, position125)
			}
			return true
		l124:
			position, tokenIndex = position124, tokenIndex124
			return false
		},
		/* 27 Spacing <- <((&('\n') '\n') | (&('\t') '\t') | (&(' ') ' '))*> */
		func() bool {
			{
				position131 := position
			l132:
				{
					position133, tokenIndex133 := position, tokenIndex
					{
						switch buffer[position] {
						case '\n':
							if buffer[position] != rune('\n') {
								goto l133
							}
							position++
							break
						case '\t':
							if buffer[position] != rune('\t') {
								goto l133
							}
							position++
							break
						default:
							if buffer[position] != rune(' ') {
								goto l133
							}
							position++
							break
						}
					}

					goto l132
				l133:
					position, tokenIndex = position133, tokenIndex133
				}
				add(ruleSpacing, position131)
			}
			return true
		},
		/* 29 Action0 <- <{ p.AddSelect() }> */
		nil,
		/* 30 Action1 <- <{ p.AddJoin() }> */
		nil,
		nil,
		/* 32 Action2 <- <{ p.SetTableName(buffer[begin:end]) }> */
		nil,
		/* 33 Action3 <- <{ p.AddJoinRow() }> */
		nil,
		/* 34 Action4 <- <{ p.SetJoinRowKey(buffer[begin:end]) }> */
		nil,
		/* 35 Action5 <- <{ p.SetJoinKey(buffer[begin:end]) }> */
		nil,
		/* 36 Action6 <- <{ p.SetJoinValue(buffer[begin:end]) }> */
		nil,
		/* 37 Action7 <- <{ p.SetTableName(buffer[begin:end]) }> */
		nil,
		/* 38 Action8 <- <{ p.SetLimit(buffer[begin:end])}> */
		nil,
		/* 39 Action9 <- <{ p.AddCryptoKey(buffer[begin:end]) }> */
		nil,
		/* 40 Action10 <- <{ p.PushWhere() }> */
		nil,
		/* 41 Action11 <- <{ p.PopWhere() }> */
		nil,
		/* 42 Action12 <- <{ p.SetWhereCommand("and") }> */
		nil,
		/* 43 Action13 <- <{ p.SetWhereCommand("or") }> */
		nil,
		/* 44 Action14 <- <{ p.InitPredicate() }> */
		nil,
		/* 45 Action15 <- <{ p.SetPredicateCommand(buffer[begin:end]) }> */
		nil,
		/* 46 Action16 <- <{ p.UsePredicateRowKey() }> */
		nil,
		/* 47 Action17 <- <{ p.AddPredicateKey(buffer[begin:end]) }> */
		nil,
		/* 48 Action18 <- <{ p.AddPredicateLiteral(buffer[begin:end])}> */
		nil,
	}
	p.rules = _rules
}
