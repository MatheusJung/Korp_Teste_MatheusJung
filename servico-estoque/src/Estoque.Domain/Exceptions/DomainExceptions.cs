namespace Estoque.Domain.Exceptions;

public class DomainException : Exception
{
    public DomainException(string message) : base(message) { }
}

public class SaldoInsuficienteException : DomainException
{
    public Guid ProdutoId { get; }
    public string CodigoProduto { get; }
    public decimal SaldoAtual { get; }
    public decimal QuantidadeSolicitada { get; }

    public SaldoInsuficienteException(
        Guid produtoId,
        string codigoProduto,
        decimal saldoAtual,
        decimal quantidadeSolicitada)
        : base($"Saldo insuficiente para o produto '{codigoProduto}'. Saldo atual: {saldoAtual}, solicitado: {quantidadeSolicitada}.")
    {
        ProdutoId = produtoId;
        CodigoProduto = codigoProduto;
        SaldoAtual = saldoAtual;
        QuantidadeSolicitada = quantidadeSolicitada;
    }
}

public class ProdutoNaoEncontradoException : DomainException
{
    public ProdutoNaoEncontradoException(Guid id)
        : base($"Produto com id '{id}' não encontrado.") { }

    public ProdutoNaoEncontradoException(string codigo)
        : base($"Produto com código '{codigo}' não encontrado.") { }
}
