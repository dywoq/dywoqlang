namespace DywoqLang.Lib.Scanner;

/// <summary>
/// Represents the tokenizer of character.
/// </summary>
/// <param name="context">The scanner context.</param>
/// <param name="character">The current character, may be null.</param>
/// <returns>The token and the boolean to tell the scanner whether to try other tokenizer for the character.</returns>
/// <exception cref="ArgumentNullException">
/// <exception cref="System.Data.NoNullAllowedException">
public delegate (Token.Token?, bool) Tokenizer(IScannerContext context, char? character);
