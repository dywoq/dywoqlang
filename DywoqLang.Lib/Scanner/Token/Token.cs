namespace DywoqLang.Lib.Scanner.Token;

public record class Token(string Literal, TokenKind Kind, TokenPosition Position);
