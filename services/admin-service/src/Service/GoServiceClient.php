<?php

namespace App\Service;

use Symfony\Contracts\HttpClient\HttpClientInterface;
use Psr\Log\LoggerInterface;

class GoServiceClient
{
    private array $serviceUrls;

    public function __construct(
        private HttpClientInterface $httpClient,
        private LoggerInterface $logger,
        string $videoServiceUrl,
        string $userServiceUrl,
        string $commentServiceUrl,
        string $historyServiceUrl
    ) {
        $this->serviceUrls = [
            'video' => $videoServiceUrl,
            'user' => $userServiceUrl,
            'comment' => $commentServiceUrl,
            'history' => $historyServiceUrl,
        ];
    }

    public function get(string $path, string $service): ?array
    {
        if (!isset($this->serviceUrls[$service])) {
            $this->logger->error("Unknown service: {$service}");
            return null;
        }

        $url = rtrim($this->serviceUrls[$service], '/') . '/' . ltrim($path, '/');

        try {
            $response = $this->httpClient->request('GET', $url, [
                'timeout' => 5,
                'headers' => [
                    'Accept' => 'application/json',
                ],
            ]);

            if ($response->getStatusCode() === 200) {
                return $response->toArray();
            }

            $this->logger->warning("Go service returned status {$response->getStatusCode()}", [
                'service' => $service,
                'url' => $url,
            ]);

            return null;
        } catch (\Exception $e) {
            $this->logger->error("Failed to call Go service: {$e->getMessage()}", [
                'service' => $service,
                'url' => $url,
                'exception' => $e,
            ]);
            return null;
        }
    }

    public function post(string $path, string $service, array $data): ?array
    {
        if (!isset($this->serviceUrls[$service])) {
            $this->logger->error("Unknown service: {$service}");
            return null;
        }

        $url = rtrim($this->serviceUrls[$service], '/') . '/' . ltrim($path, '/');

        try {
            $response = $this->httpClient->request('POST', $url, [
                'timeout' => 5,
                'json' => $data,
                'headers' => [
                    'Content-Type' => 'application/json',
                    'Accept' => 'application/json',
                ],
            ]);

            if (in_array($response->getStatusCode(), [200, 201])) {
                return $response->toArray();
            }

            return null;
        } catch (\Exception $e) {
            $this->logger->error("Failed to call Go service: {$e->getMessage()}", [
                'service' => $service,
                'url' => $url,
                'exception' => $e,
            ]);
            return null;
        }
    }

    public function put(string $path, string $service, array $data): ?array
    {
        if (!isset($this->serviceUrls[$service])) {
            $this->logger->error("Unknown service: {$service}");
            return null;
        }

        $url = rtrim($this->serviceUrls[$service], '/') . '/' . ltrim($path, '/');

        try {
            $response = $this->httpClient->request('PUT', $url, [
                'timeout' => 5,
                'json' => $data,
                'headers' => [
                    'Content-Type' => 'application/json',
                    'Accept' => 'application/json',
                ],
            ]);

            if ($response->getStatusCode() === 200) {
                return $response->toArray();
            }

            return null;
        } catch (\Exception $e) {
            $this->logger->error("Failed to call Go service: {$e->getMessage()}", [
                'service' => $service,
                'url' => $url,
                'exception' => $e,
            ]);
            return null;
        }
    }

    public function delete(string $path, string $service): bool
    {
        if (!isset($this->serviceUrls[$service])) {
            $this->logger->error("Unknown service: {$service}");
            return false;
        }

        $url = rtrim($this->serviceUrls[$service], '/') . '/' . ltrim($path, '/');

        try {
            $response = $this->httpClient->request('DELETE', $url, [
                'timeout' => 5,
            ]);

            return in_array($response->getStatusCode(), [200, 204]);
        } catch (\Exception $e) {
            $this->logger->error("Failed to call Go service: {$e->getMessage()}", [
                'service' => $service,
                'url' => $url,
                'exception' => $e,
            ]);
            return false;
        }
    }
}
