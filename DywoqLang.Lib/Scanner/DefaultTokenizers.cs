using System.Data;

namespace DywoqLang.Lib.Scanner;

public class DefaultTokenizers
{
	public static (Token.Token?, bool) Number(IScannerContext context, char? character)
	{
		if (character == null || !char.IsNumber(character.Value))
		{
			return (null, true);
		}
		var startPos = context.Position;
		context.Advance(1);
		while (char.IsNumber(context.Current ?? throw new NoNullAllowedException("Can't get current character.")))
		{
			context.Advance(1);
		}

		if (context.Peek == null || context.Peek != '.')
		{
			var subString = context.Slice(startPos, context.Position) ?? throw new NoNullAllowedException("The sliced substring cannot be null.");
			return (context.NewToken(subString, Token.TokenKind.Integer), false);
		}

		context.Advance(2);

		while (char.IsNumber(context.Current ?? throw new NoNullAllowedException("Can't get the current character")))
		{
			context.Advance(1);
		}

		var substring = context.Slice(startPos, context.Position) ?? throw new NoNullAllowedException("The sliced substring cannot be null");
		return (context.NewToken(substring, Token.TokenKind.Float), false);
	}
}