namespace DywoqLang.Lib.Scanner;

public interface IScannerReader
{
	/// <summary>
	/// Returns the future character.
	/// If there's no future character, the function returns null. 
	/// </summary>
	public char? Peek { get; }

	/// <summary>
	/// Returns the current character, else,
	/// the function returns null if the current scanner position is out of bounds of the input.
	/// </summary>
	/// <returns></returns>
	public char? Current { get; }
}
