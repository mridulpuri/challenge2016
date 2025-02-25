# Solution for Real Image Challenge 2016

## Approach

This solution implements a hierarchical permission system for movie distributors. In this system:

1. **Data Loading**  
   - We use `cities.csv` to validate region existence (cities, states, countries).  
   - We use `distributors.csv` to dynamically build distributor permissions and hierarchy.

2. **Normalization**  
   - All inputs (distributor names, regions, etc.) are normalized to uppercase and stripped of spaces.  
   - This ensures consistent matching, avoiding case and spacing issues.

3. **Permission Checking**  
   - A distributor's permissions are a subset of its parent's.  
   - `EXCLUDE` rules override `INCLUDE` rules.  
   - If a parent distributor excludes a region, all child distributors inherit that exclusion.

4. **Authorization Logic**  
   - We recursively check the distributor's parent permissions.  
   - If the parent denies permission, the child is automatically denied.  
   - If the region is explicitly excluded, deny permission.  
   - If the region is included, allow permission.  
   - Otherwise, deny.

5. **Detailed Reasoning**  
   - The solution provides reasons for why a distributor is authorized or not.

---

## Files

- **`main.go`**: Main Go program that:
  1. Loads CSV data
  2. Accepts user input (distributor name + region)
  3. Prints out authorization results

- **`cities.csv`**: Contains the canonical list of city/state/country data.

- **`distributors.csv`**: Lists each distributor's `INCLUDE`/`EXCLUDE` permissions and (optionally) a parent distributor.

- **`SOLUTION.md`**: This file, explaining the approach in detail.

---

## Running the Program

1. **Compile and Run**
   ```bash
    go run main.go
   ```

2. **When prompted**, enter:
- **Distributor Name** (e.g., `DISTRIBUTOR1`)
- **Region** in `CITY-STATE-COUNTRY` format (e.g., `CHICAGO-ILLINOIS-UNITED STATES`)


## Expected Output

The program will respond with either:
```
    DISTRIBUTOR1 is authorized to distribute in CHICAGO-ILLINOIS-UNITED STATES. Reason: Explicitly included by distributor: DISTRIBUTOR1
```
or
```
    DISTRIBUTOR1 is NOT authorized to distribute in CHICAGO-ILLINOIS-UNITED STATES. Reason: Explicitly excluded by distributor: DISTRIBUTOR1
```

---

## Example
Assuming the CSV data includes entries for `CHICAGO-ILLINOIS-UNITED STATES`:

1. **Input**:
    ```
        Enter distributor name: DISTRIBUTOR1
        Enter region (format CITY-STATE-COUNTRY): CHICAGO-ILLINOIS-UNITED STATES
    ```
2. **Output**:
    ```
    DISTRIBUTOR1 is authorized to distribute in CHICAGO-ILLINOIS-UNITED STATES. Reason: Explicitly included by distributor: DISTRIBUTOR1
    ```

---

## Enhancements and Additional Features
- **Detailed Reason** for authorization or denial.
- **Hierarchy**: The solution recursively checks parent distributors.
- **Dynamic Data**: No hardcoded distributor or region data.
- **Normalization**: Minimizes errors due to case or spacing differences.
- **Extended** to handle sub-distributors, ensuring they can't exceed their parent's permissions.

---

## Future Improvements
- **Caching**: Use caching for repeated region checks.
- **Advanced CLI**: Provide interactive commands to list distributors, visualize hierarchies, etc.
- **Web API**: Expose the authorization logic over HTTP for an external service integration.