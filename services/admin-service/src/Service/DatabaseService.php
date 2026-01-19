<?php

namespace App\Service;

use PDO;
use PDOException;

/**
 * Database Service
 * Handles all database operations for the admin service
 */
class DatabaseService
{
    private ?PDO $pdo = null;
    
    public function __construct()
    {
        $this->connect();
    }
    
    /**
     * Connect to database
     */
    private function connect(): void
    {
        $dbUrl = $_ENV['DATABASE_URL'] ?? 'postgresql://postgres:postgres@admin-db:5432/admin_service_db';
        
        // Parse PostgreSQL URL
        $parts = parse_url($dbUrl);
        $host = $parts['host'] ?? 'admin-db';
        $port = $parts['port'] ?? 5432;
        $dbname = ltrim($parts['path'] ?? '/admin_service_db', '/');
        $user = $parts['user'] ?? 'postgres';
        $pass = $parts['pass'] ?? 'postgres';
        
        $dsn = "pgsql:host={$host};port={$port};dbname={$dbname}";
        
        try {
            $this->pdo = new PDO($dsn, $user, $pass, [
                PDO::ATTR_ERRMODE => PDO::ERRMODE_EXCEPTION,
                PDO::ATTR_DEFAULT_FETCH_MODE => PDO::FETCH_ASSOC,
                PDO::ATTR_EMULATE_PREPARES => false
            ]);
        } catch (PDOException $e) {
            error_log("Database connection failed: " . $e->getMessage());
            // For now, allow service to run without database (will be created by Docker)
            $this->pdo = null;
        }
    }
    
    /**
     * Execute a query and return results
     */
    public function query(string $sql, array $params = []): array
    {
        if (!$this->pdo) {
            return [];
        }
        
        try {
            $stmt = $this->pdo->prepare($sql);
            $stmt->execute($params);
            return $stmt->fetchAll();
        } catch (PDOException $e) {
            error_log("Query error: " . $e->getMessage());
            return [];
        }
    }
    
    /**
     * Insert a record
     */
    public function insert(string $table, array $data): ?int
    {
        if (!$this->pdo) {
            return null;
        }
        
        $columns = array_keys($data);
        $placeholders = array_fill(0, count($columns), '?');
        
        $sql = sprintf(
            "INSERT INTO %s (%s) VALUES (%s) RETURNING id",
            $table,
            implode(', ', $columns),
            implode(', ', $placeholders)
        );
        
        try {
            $stmt = $this->pdo->prepare($sql);
            $stmt->execute(array_values($data));
            $result = $stmt->fetch();
            return $result['id'] ?? null;
        } catch (PDOException $e) {
            error_log("Insert error: " . $e->getMessage());
            return null;
        }
    }
    
    /**
     * Update records
     */
    public function update(string $table, array $data, array $where): bool
    {
        if (!$this->pdo) {
            return false;
        }
        
        $setClauses = [];
        foreach (array_keys($data) as $column) {
            $setClauses[] = "{$column} = ?";
        }
        
        $whereClauses = [];
        foreach (array_keys($where) as $column) {
            $whereClauses[] = "{$column} = ?";
        }
        
        $sql = sprintf(
            "UPDATE %s SET %s WHERE %s",
            $table,
            implode(', ', $setClauses),
            implode(' AND ', $whereClauses)
        );
        
        try {
            $stmt = $this->pdo->prepare($sql);
            $params = array_merge(array_values($data), array_values($where));
            return $stmt->execute($params);
        } catch (PDOException $e) {
            error_log("Update error: " . $e->getMessage());
            return false;
        }
    }
    
    /**
     * Delete records
     */
    public function delete(string $table, array $where): bool
    {
        if (!$this->pdo) {
            return false;
        }
        
        $whereClauses = [];
        foreach (array_keys($where) as $column) {
            $whereClauses[] = "{$column} = ?";
        }
        
        $sql = sprintf(
            "DELETE FROM %s WHERE %s",
            $table,
            implode(' AND ', $whereClauses)
        );
        
        try {
            $stmt = $this->pdo->prepare($sql);
            return $stmt->execute(array_values($where));
        } catch (PDOException $e) {
            error_log("Delete error: " . $e->getMessage());
            return false;
        }
    }
    
    /**
     * Initialize database tables
     */
    public function initializeTables(): void
    {
        if (!$this->pdo) {
            return;
        }
        
        $sql = file_get_contents(__DIR__ . '/../../config/schema.sql');
        
        try {
            $this->pdo->exec($sql);
        } catch (PDOException $e) {
            error_log("Table initialization error: " . $e->getMessage());
        }
    }
}
