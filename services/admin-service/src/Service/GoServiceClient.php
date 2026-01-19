<?php

namespace App\Service;

/**
 * Go Service Client
 * HTTP client for communicating with Go microservices
 */
class GoServiceClient
{
    /**
     * Make GET request to Go service
     */
    public function get(string $path, string $baseUrl): ?array
    {
        $url = rtrim($baseUrl, '/') . '/' . ltrim($path, '/');
        
        $ch = curl_init($url);
        curl_setopt_array($ch, [
            CURLOPT_RETURNTRANSFER => true,
            CURLOPT_TIMEOUT => 5,
            CURLOPT_HTTPHEADER => ['Accept: application/json']
        ]);
        
        $response = curl_exec($ch);
        
        if ($response === false) {
            error_log("cURL error: " . curl_error($ch));
            curl_close($ch);
            return null;
        }
        
        $httpCode = curl_getinfo($ch, CURLINFO_HTTP_CODE);
        curl_close($ch);
        
        if ($httpCode === 200 && $response) {
            return json_decode($response, true);
        }
        
        return null;
    }
    
    /**
     * Make POST request to Go service
     */
    public function post(string $path, string $baseUrl, array $data): ?array
    {
        $url = rtrim($baseUrl, '/') . '/' . ltrim($path, '/');
        
        $ch = curl_init($url);
        curl_setopt_array($ch, [
            CURLOPT_RETURNTRANSFER => true,
            CURLOPT_POST => true,
            CURLOPT_POSTFIELDS => json_encode($data),
            CURLOPT_TIMEOUT => 5,
            CURLOPT_HTTPHEADER => [
                'Content-Type: application/json',
                'Accept: application/json'
            ]
        ]);
        
        $response = curl_exec($ch);
        
        if ($response === false) {
            error_log("cURL error: " . curl_error($ch));
            curl_close($ch);
            return null;
        }
        
        $httpCode = curl_getinfo($ch, CURLINFO_HTTP_CODE);
        curl_close($ch);
        
        if (in_array($httpCode, [200, 201]) && $response) {
            return json_decode($response, true);
        }
        
        return null;
    }
    
    /**
     * Make PUT request to Go service
     */
    public function put(string $path, string $baseUrl, array $data): ?array
    {
        $url = rtrim($baseUrl, '/') . '/' . ltrim($path, '/');
        
        $ch = curl_init($url);
        curl_setopt_array($ch, [
            CURLOPT_RETURNTRANSFER => true,
            CURLOPT_CUSTOMREQUEST => 'PUT',
            CURLOPT_POSTFIELDS => json_encode($data),
            CURLOPT_TIMEOUT => 5,
            CURLOPT_HTTPHEADER => [
                'Content-Type: application/json',
                'Accept: application/json'
            ]
        ]);
        
        $response = curl_exec($ch);
        
        if ($response === false) {
            error_log("cURL error: " . curl_error($ch));
            curl_close($ch);
            return null;
        }
        
        $httpCode = curl_getinfo($ch, CURLINFO_HTTP_CODE);
        curl_close($ch);
        
        if ($httpCode === 200 && $response) {
            return json_decode($response, true);
        }
        
        return null;
    }
    
    /**
     * Make DELETE request to Go service
     */
    public function delete(string $path, string $baseUrl): bool
    {
        $url = rtrim($baseUrl, '/') . '/' . ltrim($path, '/');
        
        $ch = curl_init($url);
        curl_setopt_array($ch, [
            CURLOPT_RETURNTRANSFER => true,
            CURLOPT_CUSTOMREQUEST => 'DELETE',
            CURLOPT_TIMEOUT => 5
        ]);
        
        $response = curl_exec($ch);
        
        if ($response === false) {
            error_log("cURL error: " . curl_error($ch));
            curl_close($ch);
            return false;
        }
        
        $httpCode = curl_getinfo($ch, CURLINFO_HTTP_CODE);
        curl_close($ch);
        
        return in_array($httpCode, [200, 204]);
    }
}
