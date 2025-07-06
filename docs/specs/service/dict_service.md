# Specification: Dictionary Service Module

## 1. Overview

The Dictionary Service module is responsible for managing the system's data dictionaries. It is split into two distinct but related services: `DictTypeService` for managing dictionary categories (types) and `DictDataService` for managing the individual key-value pairs (data) within each type. This separation allows for a structured and organized approach to handling enumerable data, such as status codes, gender options, etc.

The module provides full CRUD functionality for both dictionary types and data, and it leverages a Redis caching layer for `DictDataService` to optimize the retrieval of dictionary entries, which are often frequently accessed and rarely changed.

## 2. Dependencies

-   **`anima/dal`**: The Data Access Layer for GORM and Redis interactions.
-   **`app/dto`**: Data Transfer Objects for request/response structuring.
-   **`app/model`**: GORM models (`SysDictType`, `SysDictData`).
-   **`common/types/constant`**: For system-wide constants like `NORMAL_STATUS`.
-   **`common/types/redis-key`**: Defines Redis keys for caching (`SysDictKey`).
-   **`context`**, **`encoding/json`**: Standard libraries for Redis operations and JSON handling.

## 3. `DictTypeService` - Dictionary Categories

This service manages the high-level dictionary types.

### 3.1. Data Structures (Key DTOs)

-   **`dto.SaveDictType`**: For creating and updating a dictionary type.
-   **`dto.DictTypeListRequest`**: For querying a list of dictionary types with filters.
-   **`dto.DictTypeListResponse`**: The structure for each item in the list response.
-   **`dto.DictTypeDetailResponse`**: The structure for a single, detailed response.

### 3.2. Functions and Logic

#### `CreateDictType`

**Pseudocode:**

```
FUNCTION CreateDictType(param: SaveDictType_DTO):
  // TDD: Test case for successful creation.
  // TDD: Test case for creation failure on duplicate DictType (if constrained in DB).
  INSTANTIATE a new SysDictType model from `param`.
  EXECUTE database INSERT.
  RETURN any resulting error.
END FUNCTION
```

#### `UpdateDictType`

**Pseudocode:**

```
FUNCTION UpdateDictType(param: SaveDictType_DTO):
  // TDD: Test case for successful update.
  // TDD: Test case for updating a non-existent dictId.
  // TDD: Test case to verify cache invalidation for related dict data.
  LOCATE SysDictType WHERE `dict_id` matches `param.DictId`.
  EXECUTE database UPDATE with data from `param`.
  // NOTE: Should invalidate the cache for this dictType.
  RETURN any resulting error.
END FUNCTION
```

#### `DeleteDictType`

**Pseudocode:**

```
FUNCTION DeleteDictType(dictIds: list of integers):
  // TDD: Test case for deleting single and multiple types.
  // TDD: Test case to ensure associated dict data is also handled (or deletion is blocked).
  EXECUTE database DELETE on SysDictType WHERE `dict_id` is IN `dictIds`.
  // NOTE: Should invalidate caches for all deleted dictTypes.
  RETURN any resulting error.
END FUNCTION
```

#### `GetDictTypeList`

**Pseudocode:**

```
FUNCTION GetDictTypeList(param: DictTypeListRequest_DTO, isPaging: boolean):
  // TDD: Test cases for all filters (DictName, DictType, Status, DateRange) and pagination.
  INITIALIZE query on SysDictType model.
  APPLY `LIKE` filters for `DictName` and `DictType`.
  APPLY `equals` filter for `Status`.
  APPLY `BETWEEN` filter for `create_time`.
  IF `isPaging`:
    APPLY `COUNT`, `OFFSET`, and `LIMIT`.
  END IF
  EXECUTE query and return the list and total count.
END FUNCTION
```

#### `GetDictTypeByDictId` / `GetDcitTypeByDictType`

**Pseudocode:**

```
FUNCTION GetDictTypeBy[Id|Type](identifier):
  // TDD: Test cases for finding by valid and invalid identifiers.
  INITIALIZE an empty response DTO.
  QUERY SysDictType model WHERE `dict_id` or `dict_type` matches the identifier.
  POPULATE DTO with the result.
  RETURN DTO.
END FUNCTION
```

---

## 4. `DictDataService` - Dictionary Entries

This service manages the specific key-value entries within a dictionary type.

### 4.1. Data Structures (Key DTOs)

-   **`dto.SaveDictData`**: For creating and updating a dictionary entry.
-   **`dto.DictDataListRequest`**: For querying a list of dictionary entries with filters.
-   **`dto.DictDataListResponse`**: The structure for each item in the list response.
-   **`dto.DictDataDetailResponse`**: The structure for a single, detailed response.

### 4.2. Functions and Logic

#### `CreateDictData`

**Pseudocode:**

```
FUNCTION CreateDictData(param: SaveDictData_DTO):
  // TDD: Test case for successful creation.
  // TDD: Test case to verify cache for the corresponding DictType is invalidated.
  INSTANTIATE a new SysDictData model from `param`.
  EXECUTE database INSERT.
  // NOTE: Should invalidate the cache for `param.DictType`.
  RETURN any resulting error.
END FUNCTION
```

#### `UpdateDictData`

**Pseudocode:**

```
FUNCTION UpdateDictData(param: SaveDictData_DTO):
  // TDD: Test case for successful update.
  // TDD: Test case to verify cache for the corresponding DictType is invalidated.
  LOCATE SysDictData WHERE `dict_code` matches `param.DictCode`.
  EXECUTE database UPDATE with data from `param`.
  // NOTE: Should invalidate the cache for `param.DictType`.
  RETURN any resulting error.
END FUNCTION
```

#### `DeleteDictData`

**Pseudocode:**

```
FUNCTION DeleteDictData(dictCodes: list of integers):
  // TDD: Test case for deleting single and multiple entries.
  // TDD: Test case to verify cache for the corresponding DictType is invalidated.
  // First, need to find the DictType of the items being deleted to clear the cache.
  EXECUTE database DELETE on SysDictData WHERE `dict_code` is IN `dictCodes`.
  // NOTE: Should invalidate the cache for the affected DictType.
  RETURN any resulting error.
END FUNCTION
```

#### `GetDictDataList`

**Pseudocode:**

```
FUNCTION GetDictDataList(param: DictDataListRequest_DTO, isPaging: boolean):
  // TDD: Test cases for all filters (DictLabel, DictType, Status) and pagination.
  INITIALIZE query on SysDictData model.
  APPLY `LIKE` filters for `DictLabel` and `DictType`.
  APPLY `equals` filter for `Status`.
  IF `isPaging`:
    APPLY `COUNT`, `OFFSET`, and `LIMIT`.
  END IF
  EXECUTE query and return the list and total count.
END FUNCTION
```

#### `GetDictDataByDictCode`

**Pseudocode:**

```
FUNCTION GetDictDataByDictCode(dictCode: integer):
  // TDD: Test case for a valid and invalid dictCode.
  INITIALIZE an empty response DTO.
  QUERY SysDictData model WHERE `dict_code` matches `dictCode`.
  POPULATE DTO with the result.
  RETURN DTO.
END FUNCTION
```

#### `GetDictDataCacheByDictType`

This function retrieves a list of dictionary entries for a given type, using a cache-aside pattern.

**Pseudocode:**

```
FUNCTION GetDictDataCacheByDictType(dictType: string):
  // TDD: Test case for a cache hit.
  // TDD: Test case for a cache miss (DB fallback).
  // TDD: Test case for an empty result from the DB.
  // TDD: Test case for a corrupted cache entry (JSON unmarshal fails).

  // 1. Cache Look-up
  ATTEMPT to get the value for `dictType` from the Redis hash `rediskey.SysDictKey`.
  IF a value is found in the cache:
    ATTEMPT to unmarshal the JSON string into a list of `DictDataListResponse`.
    IF unmarshalling is successful:
      RETURN the list from the cache.
    END IF
  END IF

  // 2. Cache Miss - Database Fallback
  CALL `GetDictDataByDictType(dictType)` to fetch the data from the database.
  LET the result be `dbDictData`.

  // 3. Cache Population
  IF `dbDictData` is not empty:
    SERIALIZE `dbDictData` into a JSON string.
    STORE the JSON string in the Redis hash `rediskey.SysDictKey` with the key `dictType`.
  END IF

  RETURN `dbDictData`.
END FUNCTION
```

#### `GetDictDataByDictType`

This is the direct database-access counterpart to the cached function.

**Pseudocode:**

```
FUNCTION GetDictDataByDictType(dictType: string):
  // TDD: Test case for a valid dictType with active entries.
  // TDD: Test case for a dictType with no entries.
  INITIALIZE an empty list of `DictDataListResponse`.
  QUERY SysDictData model WHERE `status` is `NORMAL_STATUS` AND `dict_type` matches `dictType`.
  POPULATE the list with all results.
  RETURN the list.
END FUNCTION
