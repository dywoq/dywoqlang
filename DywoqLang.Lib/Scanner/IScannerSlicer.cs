namespace DywoqLang.Lib.Scanner;

public interface IScannerSlicer
{
	/// <summary>
	/// Returns the slice of the input string with the given start and end.
	/// </summary>
	/// <param name="start">The given start.</param>
	/// <param name="end">The given end.</param>
	public string Slice(int start, int end);
}
