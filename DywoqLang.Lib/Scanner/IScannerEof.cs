namespace DywoqLang.Lib.Scanner;

public interface IScannerEof
{
	/// <summary>
	/// Reports whether the scanner reached Eof.
	/// </summary>
	public bool ReachedEof();
}
