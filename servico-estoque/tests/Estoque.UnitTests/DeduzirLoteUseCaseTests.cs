using Estoque.Application.DTOs;
using Estoque.Application.Interfaces;
using Estoque.Application.UseCases.Movimentacoes;
using Estoque.Domain.Entities;
using FluentAssertions;
using Moq;

namespace Estoque.UnitTests;

public class DeduzirLoteUseCaseTests
{
    private readonly Mock<IProdutoRepository> _repoMock = new();
    private readonly Mock<IUnitOfWork> _uowMock = new();

    private DeduzirLoteUseCase CriarUseCase() =>
        new(_repoMock.Object, _uowMock.Object);

    [Fact]
    public async Task ExecutarAsync_TodosComSaldo_DeveRetornarSucesso()
    {
        var p1 = Produto.Criar("P1", "Produto 1", 10);
        var p2 = Produto.Criar("P2", "Produto 2", 20);

        _repoMock.Setup(r => r.ObterPorIdAsync(p1.Id, default)).ReturnsAsync(p1);
        _repoMock.Setup(r => r.ObterPorIdAsync(p2.Id, default)).ReturnsAsync(p2);
        _uowMock.Setup(u => u.CommitAsync(default)).ReturnsAsync(1);

        var request = new DeduzirLoteRequest(Guid.NewGuid(), [
            new(p1.Id, 3),
            new(p2.Id, 5)
        ]);

        var resultado = await CriarUseCase().ExecutarAsync(request);

        resultado.Sucesso.Should().BeTrue();
        resultado.ItensDeduzidos.Should().HaveCount(2);
        p1.Saldo.Should().Be(7);
        p2.Saldo.Should().Be(15);
    }

    [Fact]
    public async Task ExecutarAsync_SegundoItemSemSaldo_DeveEstornarPrimeiroERetornarFalha()
    {
        var notaId = Guid.NewGuid();
        var p1 = Produto.Criar("P1", "Produto 1", 10);
        var p2 = Produto.Criar("P2", "Produto 2", 1); // saldo insuficiente

        _repoMock.Setup(r => r.ObterPorIdAsync(p1.Id, default)).ReturnsAsync(p1);
        _repoMock.Setup(r => r.ObterPorIdAsync(p2.Id, default)).ReturnsAsync(p2);
        _uowMock.Setup(u => u.CommitAsync(default)).ReturnsAsync(1);

        var request = new DeduzirLoteRequest(notaId, [
            new(p1.Id, 5),
            new(p2.Id, 10) // vai falhar
        ]);

        var resultado = await CriarUseCase().ExecutarAsync(request);

        resultado.Sucesso.Should().BeFalse();
        resultado.ProdutoFalhou.Should().Be(p2.Id);
        resultado.Erro.Should().Contain("Saldo insuficiente");

        // p1 deve ter sido estornado — saldo volta a 10
        p1.Saldo.Should().Be(10);
        // Dois commits: 1 para a dedução de p1, 1 para o estorno
        _uowMock.Verify(u => u.CommitAsync(default), Times.Exactly(2));
    }

    [Fact]
    public async Task ExecutarAsync_ProdutoNaoEncontrado_DeveRetornarFalha()
    {
        var idInexistente = Guid.NewGuid();
        _repoMock.Setup(r => r.ObterPorIdAsync(idInexistente, default)).ReturnsAsync((Produto?)null);

        var request = new DeduzirLoteRequest(Guid.NewGuid(), [new(idInexistente, 1)]);

        var resultado = await CriarUseCase().ExecutarAsync(request);

        resultado.Sucesso.Should().BeFalse();
        resultado.ProdutoFalhou.Should().Be(idInexistente);
    }
}
