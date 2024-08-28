package main

import (
	"fmt"
	"main/lexer"
	"os"
	"strconv"
)

// Order Of Presidence
// AssigmentExpr
// ObjectExpr
// BooleanExpr
// AdditiveExpr
// MultiplicitaveExpr
// CallExpr
// MemberExpr
// PrimaryExpr

type Parser struct {
	tokens            []lexer.Token
	currentTokenIndex uint
}

func (p *Parser) at() lexer.Token {
	return p.tokens[p.currentTokenIndex]
}
func (p *Parser) eat() lexer.Token {
	token := p.at()
	p.currentTokenIndex++
	return token
}
func (p *Parser) isTokenType(types ...lexer.TokenType) bool {
	currentToken := p.at()

	for _, t := range types {
		if currentToken.TokenType == t {
			return true
		}
	}
	return false
}
func (p *Parser) expect(tType lexer.TokenType) lexer.Token {
	token := p.eat()
	if token.TokenType != tType {
		fmt.Printf("Parser Error:\n Expecting: %v found: %v. Line:%v \n", tType, token.TokenType, token.Line)
		os.Exit(1)
	}
	return token
}
func (p *Parser) parseStmt() Stmt {
	if p.isTokenType(lexer.Let, lexer.Const) {
		return p.parseVarDeclaration()
	} else if p.isTokenType(lexer.Fn) {
		return p.parseFnDecralation()
	} else if p.isTokenType(lexer.If) {
		return p.parseIfStmt()
	}

	return p.parseExpr()
}
func (p *Parser) parseVarDeclaration() VarDeclaration {
	isConstant := p.eat().TokenType == lexer.Const
	identifier := p.expect(lexer.Identifier).Value

	if p.isTokenType(lexer.Semicolon) {
		p.eat()
		if isConstant {
			println("Cannot initialize constant variable without value")
			os.Exit(1)
		}
		return VarDeclaration{constant: false, identifier: identifier}
	}
	p.expect(lexer.Equals)
	decralation := VarDeclaration{
		constant: isConstant, identifier: identifier, value: p.parseExpr(),
	}
	nextToken := p.eat()
	if nextToken.TokenType != lexer.Semicolon {
		fmt.Printf("Missing semicolon at the end of variable declaration: %v Line:%v\n", decralation.identifier, nextToken.Line)
		os.Exit(1)
	}

	return decralation

}
func (p *Parser) parseFnDecralation() FunctionDeclaration {
	p.eat()
	name := p.expect(lexer.Identifier).Value
	args := p.parseArgs()
	params := make([]string, len(args))

	for i, arg := range args {
		v, ok := arg.(Identifier)
		if !ok {
			fmt.Println("Expect strings as parametrs inside function declaration")
		}
		params[i] = v.symbol
	}
	p.expect(lexer.OpenBrace)

	body := make([]Stmt, 0)
	for !p.isTokenType(lexer.EOF, lexer.CloseBrace) {
		body = append(body, p.parseStmt())
	}
	p.expect(lexer.CloseBrace)
	return FunctionDeclaration{name: name, parameters: params, body: body}
}
func (p *Parser) parseIfStmt() IfStmt {
	p.eat()

	p.expect(lexer.OpenParen)
	condition := p.parseExpr()
	p.expect(lexer.CloseParen)

	p.expect(lexer.OpenBrace)
	body := make([]Stmt, 0)
	for !p.isTokenType(lexer.EOF, lexer.CloseBrace) {
		body = append(body, p.parseStmt())
	}
	p.expect(lexer.CloseBrace)

	alternative := make([]Stmt, 0)
	if p.isTokenType(lexer.Else) {
		p.eat()
		p.expect(lexer.OpenBrace)
		for !p.isTokenType(lexer.EOF, lexer.CloseBrace) {
			alternative = append(alternative, p.parseStmt())
		}
		p.expect(lexer.CloseBrace)
	}

	return IfStmt{condition, body, alternative}
}
func (p *Parser) parseExpr() Expr {

	return p.parseAssignmentExpr()
}
func (p *Parser) parseAssignmentExpr() Expr {
	left := p.parseObjectExpr()

	if p.isTokenType(lexer.Equals) {
		p.eat()
		value := p.parseAssignmentExpr()
		return AssigmentExpr{value: value, assigne: left}
	}
	return left
}
func (p *Parser) parseObjectExpr() Expr {
	if !p.isTokenType(lexer.OpenBrace) {
		return p.parseBooleanExpr()
	}
	p.eat()
	properties := make([]Property, 0)

	for !p.isTokenType(lexer.EOF, lexer.CloseBrace) {
		key := p.expect(lexer.Identifier).Value

		if p.isTokenType(lexer.Coma) {
			p.eat()
			properties = append(properties, Property{key: key, value: nil})
			continue
		} else if p.isTokenType(lexer.CloseBrace) {
			properties = append(properties, Property{key: key, value: nil})
			continue
		}

		p.expect(lexer.Colon)
		value := p.parseExpr()

		properties = append(properties, Property{key, value})

		if !p.isTokenType(lexer.CloseBrace) {
			p.expect(lexer.Coma)
		}
	}

	p.expect(lexer.CloseBrace)

	return ObjectLiteral{properties}
}
func (p *Parser) parseBooleanExpr() Expr {
	left := p.parseUnaryExpr()

	for p.isTokenType(lexer.EqualsEquals, lexer.NotEquals, lexer.LessThan, lexer.GreaterThan, lexer.LessThanOrEquals, lexer.GreaterThanOrEquals) {
		operator := p.eat().Value
		right := p.parseUnaryExpr()
		left = BooleanExpr{left, right, operator}
	}
	return left
}
func (p *Parser) parseUnaryExpr() Expr {

	if p.isTokenType(lexer.Not) {
		p.eat()
		return UnaryExpression{operator: "!", operand: p.parseAdditiveExpr()}
	}

	return p.parseAdditiveExpr()
}
func (p *Parser) parseAdditiveExpr() Expr {
	left := p.parseMultiplicitaveExpr()

	for p.at().Value == "+" || p.at().Value == "-" {
		operator := p.eat().Value
		right := p.parseMultiplicitaveExpr()
		left = BinaryExpr{left, right, operator}
	}

	return left
}
func (p *Parser) parseMultiplicitaveExpr() Expr {
	left := p.parseCallMemberExpr()

	for p.at().Value == "/" || p.at().Value == "*" || p.at().Value == "%" {
		operator := p.eat().Value
		right := p.parseCallMemberExpr()
		left = BinaryExpr{left, right, operator}

	}
	return left
}
func (p *Parser) parseCallMemberExpr() Expr {
	member := p.parseMemberExpr()

	if p.isTokenType(lexer.OpenParen) {
		return p.parseCallExpr(member)
	}
	return member
}
func (p *Parser) parseCallExpr(caller Expr) Expr {
	var callExpr Expr = CallExpr{caller: caller, args: p.parseArgs()}

	if p.isTokenType(lexer.OpenParen) {
		callExpr = p.parseCallExpr(callExpr)
	}
	return callExpr
}
func (p *Parser) parseArgs() []Expr {
	p.expect(lexer.OpenParen)
	args := make([]Expr, 0)

	if p.isTokenType(lexer.CloseParen) {
		p.eat()
		return args
	} else {
		args = p.parseArgsList()
		p.expect(lexer.CloseParen)
	}
	return args
}
func (p *Parser) parseArgsList() []Expr {
	args := []Expr{p.parseAssignmentExpr()}

	for p.isTokenType(lexer.Coma) && !p.isTokenType(lexer.EOF) {
		p.eat()
		args = append(args, p.parseAssignmentExpr())
	}

	return args
}
func (p *Parser) parseMemberExpr() Expr {
	object := p.parsePrimaryExpr()

	for p.isTokenType(lexer.Dot, lexer.OpenBracket) {
		operator := p.eat()
		var property Expr
		var computed bool
		if operator.TokenType == lexer.Dot {
			computed = false
			property = p.parsePrimaryExpr()
			if _, ok := property.(Identifier); !ok {
				fmt.Println("Cannot use dot operator without right hand side being an identifier")
			}
		} else {
			computed = true
			property = p.parseExpr()
			p.expect(lexer.CloseBracket)
		}
		object = MemberExpr{object, property, computed}
	}
	return object
}
func (p *Parser) parsePrimaryExpr() Expr {
	switch p.at().TokenType {
	case lexer.Identifier:
		return Identifier{symbol: p.eat().Value}
	case lexer.Number:
		value, err := strconv.ParseInt(p.at().Value, 10, 64)
		if err != nil {
			fmt.Printf("Error while parsing number literal: '%v' \n", p.eat().Value)
			os.Exit(1)
		}
		p.eat()
		return NumericLiteral{value}
	case lexer.String:
		return StringLiteral{value: p.eat().Value}
	case lexer.OpenParen:
		p.eat()
		value := p.parseExpr()
		if p.eat().TokenType != lexer.CloseParen {
			fmt.Printf("Parser Error:\n - Expecting: CloseParen\n")
			os.Exit(1)

		}
		return value
	case lexer.OpenBracket:
		p.eat()
		elements := p.parseArgsList()
		p.expect(lexer.CloseBracket)
		return ArrayLiteral{elements}
	default:
		token := p.eat()
		fmt.Printf("Unexpected token found during parsing: '%v' Line:%v\n", token.Value, token.Line)
		os.Exit(1)
		panic("Unreachable code")
	}
}

func produceAst(sourceCode string) Program {
	tokens := lexer.Tokenize(sourceCode)

	parser := Parser{tokens: tokens, currentTokenIndex: 0}
	program := Program{body: make([]Stmt, 0)}
	for parser.at().TokenType != lexer.EOF {
		program.body = append(program.body, parser.parseStmt())
	}

	return program
}
