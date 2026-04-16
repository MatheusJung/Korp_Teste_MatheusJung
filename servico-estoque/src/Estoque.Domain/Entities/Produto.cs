namespace Estoque.Domain.Entities;

public class Produto
{
    public Guid Id { get; private set; }
    public string Codigo { get; private set; } = string.Empty;
    public string Descricao { get; private set; } = string.Empty;
    public decimal Saldo { get; private set; }
    public DateTime CriadoEm { get; private set; }
    public DateTime AtualizadoEm { get; private set; }

    // RowVersion para controle de concorrência otimista no EF Core
    public byte[] RowVersion { get; set; } = [];

    private readonly List<Movimentacao> _movimentacoes = [];
    public IReadOnlyCollection<Movimentacao> Movimentacoes => _movimentacoes.AsReadOnly();

    private Produto() { }

    public static Produto Criar(string codigo, string descricao, decimal saldoInicial = 0)
    {
        if (string.IsNullOrWhiteSpace(codigo))
            throw new Domain.Exceptions.DomainException("Código do produto é obrigatório.");

        if (string.IsNullOrWhiteSpace(descricao))
            throw new Domain.Exceptions.DomainException("Descrição do produto é obrigatória.");

        if (saldoInicial < 0)
            throw new Domain.Exceptions.DomainException("Saldo inicial não pode ser negativo.");

        var produto = new Produto
        {
            Id = Guid.NewGuid(),
            Codigo = codigo.Trim().ToUpper(),
            Descricao = descricao.Trim(),
            Saldo = saldoInicial,
            CriadoEm = DateTime.UtcNow,
            AtualizadoEm = DateTime.UtcNow
        };

        // Registra criação
        produto._movimentacoes.Add(Movimentacao.CriarCriacao(produto.Id, saldoInicial));

        return produto;
    }

    public void Deduzir(decimal quantidade, Guid? notaFiscalId = null)
    {
        if (quantidade <= 0)
            throw new Domain.Exceptions.DomainException("Quantidade para dedução deve ser positiva.");

        if (Saldo < quantidade)
            throw new Domain.Exceptions.SaldoInsuficienteException(Id, Codigo, Saldo, quantidade);

        var saldoAnterior = Saldo;
        Saldo -= quantidade;
        AtualizadoEm = DateTime.UtcNow;

        _movimentacoes.Add(Movimentacao.CriarSaida(Id, quantidade, saldoAnterior, Saldo, notaFiscalId));
    }

    public void Estornar(decimal quantidade, Guid? notaFiscalId = null)
    {
        if (quantidade <= 0)
            throw new Domain.Exceptions.DomainException("Quantidade para estorno deve ser positiva.");

        var saldoAnterior = Saldo;
        Saldo += quantidade;
        AtualizadoEm = DateTime.UtcNow;

        _movimentacoes.Add(Movimentacao.CriarEntrada(Id, quantidade, saldoAnterior, Saldo, notaFiscalId, isEstorno: true));
    }

    public void AdicionarEntrada(decimal quantidade, Guid? notaFiscalId = null)
    {
        if (quantidade <= 0)
            throw new Domain.Exceptions.DomainException("Quantidade de entrada deve ser positiva.");

        var saldoAnterior = Saldo;
        Saldo += quantidade;
        AtualizadoEm = DateTime.UtcNow;

        _movimentacoes.Add(Movimentacao.CriarEntrada(Id, quantidade, saldoAnterior, Saldo, notaFiscalId));
    }
}
