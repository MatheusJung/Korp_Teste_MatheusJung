using Microsoft.EntityFrameworkCore;
using Microsoft.EntityFrameworkCore.Design;

namespace BillingService.Infrastructure.Data
{
    public class BillingDbContextFactory : IDesignTimeDbContextFactory<BillingDbContext>
    {
        public BillingDbContext CreateDbContext(string[] args)
        {
            var optionsBuilder = new DbContextOptionsBuilder<BillingDbContext>();

            optionsBuilder.UseSqlServer(
                "Server=localhost,1433;Database=billingdb;User Id=sa;Password=YourStrong@Password1;TrustServerCertificate=True;",
                sqlOptions => sqlOptions.EnableRetryOnFailure() 
            );
            return new BillingDbContext(optionsBuilder.Options);
        }
    }
}
