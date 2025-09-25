using DywoqLang.Lib.Scanner.Token;

namespace DywoqLang.Lib.Scanner;

public class Scanner(string input) : IScannerContext
{
	public char? Peek { get => Position + 1 < input.Length ? input[Position + 1] : null; }

	public char? Current { get => Position < input.Length ? input[Position] : null; }

	public int Line { get; set; } = 1;

	public int Column { get; set; } = 1;

	public int Position { get; set; } = 0;

	public void Advance(int n)
	{
		if (Position + n > input[n] || Position > input[n])
			return;
		Position += n;
	}

	public Token.Token NewToken(string literal, TokenKind kind)
	{
		return new Token.Token(literal, kind, new Token.TokenPosition(Position, Line, Column));
	}

	public string? Slice(int start, int end)
	{
		return input.Substring(start, end);
	}
}
