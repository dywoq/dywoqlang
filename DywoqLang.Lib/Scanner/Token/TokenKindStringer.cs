using System.Collections.ObjectModel;

namespace DywoqLang.Lib.Scanner.Token;

public class TokenKindStringer
{
	/// <summary>
	/// Map with the string versions of enumeration <see cref="TokenKind"></see>.
	/// </summary>
	public readonly static ReadOnlyDictionary<TokenKind, string> Map = new(
		new Dictionary<TokenKind, string>()
		{
			{TokenKind.Identifier, "identifier"},
			{TokenKind.BaseInstruction, "base instruction"},
			{TokenKind.Float, "float"},
			{TokenKind.Integer, "integer"},
			{TokenKind.Registry, "registry"},
			{TokenKind.String, "string"},
			{TokenKind.Type,  "type"}
		}
	);

	/// <summary>
	/// Converts TokenKind enumerations to string.
	/// </summary>
	public static string String(TokenKind kind)
	{
		return Map[kind];
	}
}
