namespace DywoqLang.Lib.Scanner;

public interface IScannerPosition
{
	/// <summary>
	/// Gets the current line.
	/// </summary>
	public int Line { get; set; }

	/// <summary>
	/// Gets the current column.
	/// </summary>
	public int Column { get; set; }

	/// <summary>
	/// Gets the current position.
	/// </summary>
	public int Position { get; set; }
}
