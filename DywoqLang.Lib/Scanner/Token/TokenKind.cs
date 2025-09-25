namespace DywoqLang.Lib.Scanner.Token;

/// <summary>
/// Represents the token kind.
/// </summary>
public enum TokenKind
{
	Identifier,
	BaseInstruction,
	Integer,
	Float,
	String,
	Registry,
	Type,
	Illegal
}
