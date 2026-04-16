using Estoque.Domain.Entities;
using Estoque.Domain.Enums;
using Microsoft.EntityFrameworkCore;
using Microsoft.EntityFrameworkCore.Metadata.Builders;

namespace Estoque.Infrastructure.Persistence.Configurations;

public class ProdutoConfiguration : IEntityTypeConfiguration<Produto>
{
    public void Configure(EntityTypeBuilder<Produto> builder)
    {
        builder.ToTable("Produtos");

        builder.HasKey(p => p.Id);

        builder.Property(p => p.Codigo)
            .IsRequired()
            .HasMaxLength(50);

        builder.Property(p => p.Descricao)
            .IsRequired()
            .HasMaxLength(200);

        builder.Property(p => p.Saldo)
            .HasColumnType("decimal(18,4)")
            .IsRequired();

        builder.Property(p => p.CriadoEm).IsRequired();
        builder.Property(p => p.AtualizadoEm).IsRequired();

        // Concorrência otimista — SQL Server atualiza automaticamente via timestamp
        builder.Property(p => p.RowVersion)
            .IsRowVersion()
            .IsConcurrencyToken();

        builder.HasIndex(p => p.Codigo).IsUnique();

        builder.HasMany(p => p.Movimentacoes)
            .WithOne()
            .HasForeignKey(m => m.ProdutoId)
            .OnDelete(DeleteBehavior.Cascade);
    }
}

public class MovimentacaoConfiguration : IEntityTypeConfiguration<Movimentacao>
{
    public void Configure(EntityTypeBuilder<Movimentacao> builder)
    {
        builder.ToTable("Movimentacoes");

        builder.HasKey(m => m.Id);

        builder.Property(m => m.Tipo)
            .HasConversion<string>()
            .IsRequired();

        builder.Property(m => m.Quantidade)
            .HasColumnType("decimal(18,4)")
            .IsRequired();

        builder.Property(m => m.SaldoAnterior)
            .HasColumnType("decimal(18,4)")
            .IsRequired();

        builder.Property(m => m.SaldoResultante)
            .HasColumnType("decimal(18,4)")
            .IsRequired();

        builder.Property(m => m.NotaFiscalId).IsRequired(false);
        builder.Property(m => m.IsEstorno).IsRequired();
        builder.Property(m => m.OcorridoEm).IsRequired();

        builder.HasIndex(m => m.ProdutoId);
        builder.HasIndex(m => m.NotaFiscalId);
        builder.HasIndex(m => m.OcorridoEm);
    }
}
