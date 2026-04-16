using Estoque.Domain.Enums;

namespace Estoque.Domain.Entities;

public class Movimentacao
{
    public Guid Id { get; private set; }
    public Guid ProdutoId { get; private set; }
    public TipoMovimentacao Tipo { get; private set; }
    public decimal Quantidade { get; private set; }
    public decimal SaldoAnterior { get; private set; }
    public decimal SaldoResultante { get; private set; }
    public Guid? NotaFiscalId { get; private set; }
    public bool IsEstorno { get; private set; }
    public DateTime OcorridoEm { get; private set; }

    private Movimentacao() { }

    public static Movimentacao CriarSaida(
        Guid produtoId,
        decimal quantidade,
        decimal saldoAnterior,
        decimal saldoResultante,
        Guid? notaFiscalId = null)
    {
        return new Movimentacao
        {
            Id = Guid.NewGuid(),
            ProdutoId = produtoId,
            Tipo = TipoMovimentacao.Saida,
            Quantidade = quantidade,
            SaldoAnterior = saldoAnterior,
            SaldoResultante = saldoResultante,
            NotaFiscalId = notaFiscalId,
            IsEstorno = false,
            OcorridoEm = DateTime.UtcNow
        };
    }

    public static Movimentacao CriarEntrada(
        Guid produtoId,
        decimal quantidade,
        decimal saldoAnterior,
        decimal saldoResultante,
        Guid? notaFiscalId = null,
        bool isEstorno = false)
    {
        return new Movimentacao
        {
            Id = Guid.NewGuid(),
            ProdutoId = produtoId,
            Tipo = TipoMovimentacao.Entrada,
            Quantidade = quantidade,
            SaldoAnterior = saldoAnterior,
            SaldoResultante = saldoResultante,
            NotaFiscalId = notaFiscalId,
            IsEstorno = isEstorno,
            OcorridoEm = DateTime.UtcNow
        };
    }

    public static Movimentacao CriarCriacao(Guid produtoId, decimal saldoInicial = 0)
    {
        return new Movimentacao
        {
            Id = Guid.NewGuid(),
            ProdutoId = produtoId,
            Tipo = TipoMovimentacao.Criacao,
            Quantidade = saldoInicial,
            SaldoAnterior = 0,
            SaldoResultante = saldoInicial,
            NotaFiscalId = null,
            IsEstorno = false,
            OcorridoEm = DateTime.UtcNow
        };
    }
}
