using Iris.Subscription.Featuretoggle.Data.Models;
using Iris.Subscription.Featuretoggle.Repository;
using Microsoft.Extensions.Logging;
using System.Diagnostics;
using System.Runtime.CompilerServices;
using System.Text.Json;
namespace Iris.Subscription.Featuretoggle.Services;

public class TenantService(ITenantRepository tenantRepository, ILogger<TenantService> logger, IPlatformCache platformCache) : ITenantService
{
    private readonly ITenantRepository _tenantRepository = tenantRepository;
    private readonly ILogger<TenantService> _logger = logger;
    private readonly IPlatformCache _platformCache = platformCache;

    private const string FeaturesForTenantCachePrefix = "FeaturesForTenant_";
    private const string FeatureForTenantCachePrefix = "FeatureForTenant_";
public async Task<GetFeaturesPerTenantResponse> GetFeaturesPerTenantAsync(
        GetFeaturesPerTenantRequest request,
        CancellationToken cancellationToken = default,
        [CallerMemberName] string callerName = "")
    {
        var cacheKey = FeaturesForTenantCachePrefix + request.TenantId;

        var featuresTenantEntityCached = await _platformCache.GetValueAsync(cacheKey);
        if (!string.IsNullOrEmpty(featuresTenantEntityCached))
        {
            _logger.LogDebug(
                    "{Caller}: '{CacheKey}' has been retrieved from cache: {Value} ",
                    callerName,
                    cacheKey,
                    featuresTenantEntityCached);

            return JsonSerializer.Deserialize<GetFeaturesPerTenantResponse>(featuresTenantEntityCached);
        }
        var featuresTenantEntity = await _tenantRepository.GetListOfFeaturesForTenantAsync(request, cancellationToken);
        await _platformCache.SetValueAsync(cacheKey, JsonSerializer.Serialize(featuresTenantEntity));
        return featuresTenantEntity;
    }

    public async Task<GetFeatureStatusForTenantResponse> GetFeatureStatusForTenantAsync(
        GetFeatureStatusForTenantRequest request,
        CancellationToken cancellationToken = default,
        [CallerMemberName] string callerName = "")
    {
        var cacheKey = $"{FeatureForTenantCachePrefix}{request.TenantId.ToString().ToLowerInvariant()}_{request.FeatureSystemName.ToString().ToLowerInvariant()}";
        var tenantFeatureCache = await _platformCache.GetValueAsync(cacheKey, cancellationToken);
        if (!string.IsNullOrEmpty(tenantFeatureCache))
        {
            _logger.LogDebug(
                    "{Caller}: '{CacheKey}' has been retrieved from cache: {Value} ",
                    callerName,
                    cacheKey,
                    tenantFeatureCache);
            return JsonSerializer.Deserialize<GetFeatureStatusForTenantResponse>(tenantFeatureCache);
        }

        var featureTenantEntity = await _tenantRepository.GetFeatureStatusForTenantAsync(request, cancellationToken);
        var value = JsonSerializer.Serialize(featureTenantEntity);
        await _platformCache.SetValueAsync(cacheKey, value, cancellationToken: cancellationToken, callerName: callerName);

        return featureTenantEntity;
    }

    public async Task<long> ClearTenantsCacheAsync()
    {
        var stopWatch = Stopwatch.StartNew();
        var cachedResult = await _platformCache.GetAllKeys($"{FeaturesForTenantCachePrefix}*");
        _logger.LogDebug("A number of {NumberOfFeatures} Tenants' features cache keys has been retrieved.", cachedResult.Length);
        var deletedCachedKeys = await _platformCache.DeleteKeys(cachedResult);
        stopWatch.Stop();
        _logger.LogDebug("Tenants' features cache has been cleared in {Time}ms.", stopWatch.Elapsed.TotalMilliseconds);
        return deletedCachedKeys;
    }

    public async Task<long> ClearTenantFeatureCacheAsync()
    {
        var stopWatch = Stopwatch.StartNew();

        var cachedResult = await _platformCache.GetAllKeys($"{FeatureForTenantCachePrefix}*");
        _logger.LogDebug("A number of {NumberOfFeatures} Tenant's features cache keys has been retrieved.", cachedResult.Length);

        var deletedCachedKeys = await _platformCache.DeleteKeys(cachedResult);
        stopWatch.Stop();
        _logger.LogDebug("Tenant's features cache has been cleared in {Time}ms.", stopWatch.Elapsed.TotalMilliseconds);

        return deletedCachedKeys;
    }
}
