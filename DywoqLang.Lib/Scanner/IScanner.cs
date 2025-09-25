namespace DywoqLang.Lib.Scanner;

public interface IScanner
{
	/// <summary>
	/// Scans the input and returns a list of tokens.
	/// </summary>
	public List<Token.Token>? Scan(string input);
}
