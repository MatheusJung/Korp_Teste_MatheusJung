using Estoque.Domain.Entities;
using Estoque.Domain.Exceptions;
using FluentAssertions;

namespace Estoque.UnitTests;

public class ProdutoTests
{
    [Fact]
    public void Criar_ComDadosValidos_DeveRetornarProduto()
    {
        var produto = Produto.Criar("PROD-001", "Produto Teste", 100);

        produto.Codigo.Should().Be("PROD-001");
        produto.Descricao.Should().Be("Produto Teste");
        produto.Saldo.Should().Be(100);
        produto.Id.Should().NotBeEmpty();
    }

    [Theory]
    [InlineData("", "Descricao")]
    [InlineData("  ", "Descricao")]
    public void Criar_ComCodigoVazio_DeveLancarDomainException(string codigo, string descricao)
    {
        var act = () => Produto.Criar(codigo, descricao);
        act.Should().Throw<DomainException>();
    }

    [Fact]
    public void Criar_ComSaldoNegativo_DeveLancarDomainException()
    {
        var act = () => Produto.Criar("COD", "Desc", -1);
        act.Should().Throw<DomainException>();
    }

    [Fact]
    public void Deduzir_ComSaldoSuficiente_DeveReduzirSaldo()
    {
        var produto = Produto.Criar("P001", "Produto", 10);

        produto.Deduzir(3);

        produto.Saldo.Should().Be(7);
        produto.Movimentacoes.Should().HaveCount(1);
        produto.Movimentacoes.First().SaldoAnterior.Should().Be(10);
        produto.Movimentacoes.First().SaldoResultante.Should().Be(7);
    }

    [Fact]
    public void Deduzir_ComSaldoInsuficiente_DeveLancarSaldoInsuficienteException()
    {
        var produto = Produto.Criar("P001", "Produto", 1);

        var act = () => produto.Deduzir(2);

        act.Should().Throw<SaldoInsuficienteException>()
            .Which.SaldoAtual.Should().Be(1);
    }

    [Fact]
    public void Estornar_DeveAumentarSaldoERegistrarMovimentacao()
    {
        var produto = Produto.Criar("P001", "Produto", 5);
        produto.Deduzir(5);

        produto.Estornar(5, Guid.NewGuid());

        produto.Saldo.Should().Be(5);
        produto.Movimentacoes.Should().HaveCount(2);
        produto.Movimentacoes.Last().IsEstorno.Should().BeTrue();
    }

    [Fact]
    public void AdicionarEntrada_DeveAumentarSaldo()
    {
        var produto = Produto.Criar("P001", "Produto", 10);

        produto.AdicionarEntrada(5);

        produto.Saldo.Should().Be(15);
    }

    [Fact]
    public void Deduzir_ExatamenteOSaldo_DeveFuncionar()
    {
        var produto = Produto.Criar("P001", "Produto", 5);

        produto.Deduzir(5);

        produto.Saldo.Should().Be(0);
    }
}
