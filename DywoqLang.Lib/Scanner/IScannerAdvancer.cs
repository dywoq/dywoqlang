namespace DywoqLang.Lib.Scanner;

public interface IScannerAdvancer
{
	/// <summary>
	/// Advances to the next position by n.
	/// If the current position + n is out of bounds of the input,
	/// or the current position is out of bounds, the function does nothing.
	/// </summary>
	public void Advance(int n);
}
