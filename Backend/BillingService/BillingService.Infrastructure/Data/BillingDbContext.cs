using BillingService.Domain.Entities;
using Microsoft.EntityFrameworkCore;

namespace BillingService.Infrastructure.Data
{
    public class BillingDbContext : DbContext
    {
        public BillingDbContext(DbContextOptions<BillingDbContext> options) : base(options) { }

        public DbSet<Invoice> Invoices => Set<Invoice>();
        public DbSet<InvoiceItem> InvoiceItems => Set<InvoiceItem>();

        protected override void OnModelCreating(ModelBuilder modelBuilder)
        {
            base.OnModelCreating(modelBuilder); // importante manter

            // Configuração da entidade Invoice
            modelBuilder.Entity<Invoice>(entity =>
            {
                entity.HasKey(i => i.Id);

                entity.Property(i => i.SequentialNumber)
                      .IsRequired();

                entity.Property(i => i.Status)
                      .IsRequired();

                // Configura a relação com InvoiceItem usando backing field _items
                entity.HasMany(i => i.Items)              // navigation property
                      .WithOne(ii => ii.Invoice)
                      .HasForeignKey(ii => ii.InvoiceId)
                      .OnDelete(DeleteBehavior.Cascade);

                entity.Navigation(i => i.Items)
                      .UsePropertyAccessMode(PropertyAccessMode.Field);

                // Índice único para SequentialNumber
                entity.HasIndex(i => i.SequentialNumber)
                      .IsUnique();
            });

            // Configuração da entidade InvoiceItem
            modelBuilder.Entity<InvoiceItem>(entity =>
            {
                entity.HasKey(ii => ii.Id);

                entity.Property(ii => ii.ProductCode)
                      .IsRequired();

                entity.Property(ii => ii.Description)
                      .IsRequired();

                entity.Property(ii => ii.Quantity)
                      .IsRequired();
            });
        }
    }
}
