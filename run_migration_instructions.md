# Supabase Migration Instructions

## Step 1: Run the Database Migration

1. Go to your Supabase dashboard: https://supabase.com/dashboard
2. Select your "ethosview" project
3. Go to the "SQL Editor" tab
4. Copy and paste the contents of `supabase_migration.sql` into the editor
5. Click "Run" to execute the migration

## Step 2: Verify the Migration

After running the migration, you should see:
- 3 tables created: `users`, `companies`, `esg_scores`
- Sample data inserted (10 companies and 50 ESG scores)
- Indexes and triggers created
- RLS policies enabled

## Step 3: Test the Connection

Once the migration is complete, run:
```bash
go run test_supabase_connection.go
```

## Migration SQL Content

The migration file contains:
- UUID-based table structure
- Sample companies (Apple, Microsoft, Amazon, etc.)
- Sample ESG scores with realistic data
- Proper indexes for performance
- Row Level Security policies

## Troubleshooting

If you get connection errors:
1. Check that your Supabase project is active
2. Verify the database password is correct
3. Ensure your IP is not blocked by Supabase
4. Check that the project URL is correct

## Next Steps After Migration

1. Test the connection
2. Start the backend server
3. Start the frontend development server
4. Verify the application works with Supabase
