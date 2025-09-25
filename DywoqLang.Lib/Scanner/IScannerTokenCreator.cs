using DywoqLang.Lib.Scanner.Token;

namespace DywoqLang.Lib.Scanner;

public interface IScannerTokenCreator
{
	/// <summary>
	/// Returns a new token with the automatically set position.
	/// </summary>
	public Token.Token NewToken(string literal, TokenKind kind);
}
