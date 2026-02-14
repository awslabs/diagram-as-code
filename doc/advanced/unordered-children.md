# UnorderedChildren Feature

## Overview

The `UnorderedChildren` feature allows automatic reordering of child resources based on link connections to minimize visual overlap and crossing. When enabled, the system analyzes link relationships and rearranges children to create more optimal diagram layouts.

## Motivation

When creating diagrams with links between resources, the order of children in the YAML file may not be optimal for visual clarity. Links can overlap intermediate resources, creating confusing diagrams.

### Example Problem

```yaml
ResourceA:
  Direction: horizontal
  Children:
    - ResourceB  # Contains Resource1, Resource2
    - ResourceC  # Contains Resource3, Resource4
    - ResourceD  # Contains Resource5, Resource6

Links:
  - Source: Resource1  # In ResourceB
    Target: Resource6  # In ResourceD
```

When drawing an orthogonal link from Resource1 to Resource6, the link path overlaps with multiple resources:

1. **Resource2** (in ResourceB): The link exits Resource1 moving right (toward LCA), overlapping Resource2 if it's positioned to the right of Resource1
2. **ResourceC and its children** (Resource3, Resource4): The link crosses through the intermediate container
3. **Resource5** (in ResourceD): The link enters ResourceD from the left, overlapping Resource5 if it's positioned to the left of Resource6

![Problem: Link overlaps with multiple resources](./unordered-children-before.png)

*Figure: Without UnorderedChildren, the link from Resource1 to Resource6 crosses over intermediate resources*

To eliminate these overlaps, we need multi-level reordering:
- **Within ResourceB**: Move Resource1 to the rightmost position (so the link exits cleanly)
- **Within ResourceD**: Move Resource6 to the leftmost position (so the link enters cleanly)
- **At ResourceA (LCA)**: Move ResourceD adjacent to ResourceB to minimize the crossing distance

## Solution

By setting `UnorderedChildren: true` on resources at multiple levels, the system automatically reorders children to minimize link overlaps:

```yaml
ResourceA:
  Direction: horizontal
  Options:
    UnorderedChildren: true
  Children:
    - ResourceB
    - ResourceC
    - ResourceD

ResourceB:
  Direction: horizontal
  Options:
    UnorderedChildren: true
  Children:
    - Resource1
    - Resource2

ResourceD:
  Direction: horizontal
  Options:
    UnorderedChildren: true
  Children:
    - Resource5
    - Resource6
```

The system will analyze the link and perform multi-level reordering:
1. **ResourceB**: Reorder to place Resource1 at the rightmost position
2. **ResourceD**: Reorder to place Resource6 at the leftmost position
3. **ResourceA**: Reorder to place ResourceD adjacent to ResourceB

![Solution: After reordering with UnorderedChildren](./unordered-children-after.png)

*Figure: With UnorderedChildren enabled, children are automatically reordered to minimize link overlaps*

### Resource Hierarchy Comparison

The following diagrams show the parent-child relationship tree and how UnorderedChildren reorders the structure:

**Before (UnorderedChildren disabled):**

![Hierarchy Tree Before](./hierarchy-tree-before.png)

*Original order: ResourceA → [B,C,D], ResourceB → [1,2], ResourceD → [5,6]. Gray lines show parent-child relationships. Red link from Resource1 to Resource6 crosses over Resources 2,3,4,5.*

**After (UnorderedChildren enabled):**

![Hierarchy Tree After](./hierarchy-tree-after.png)

*Reordered: ResourceA → [B,D,C], ResourceB → [2,1], ResourceD → [6,5]. Blue lines indicate reordered relationships. Green dashed lines show which elements were swapped. Red link is now optimized with Resource1 and Resource6 adjacent.*

## Algorithm

The reordering algorithm processes each link once and performs the following steps:

### Step 1: Find Lowest Common Ancestor (LCA)

Find the LCA between the link's source and target resources.

```
Example: Resource1 → Resource6
LCA = ResourceA
```

### Step 2: Check LCA Direction

Verify the LCA's layout direction (horizontal or vertical).

```
ResourceA.Direction = "horizontal"
```

### Step 3: Identify LCA Children

Determine which direct children of the LCA contain the source and target.

```
Source (Resource1) belongs to: ResourceB
Target (Resource6) belongs to: ResourceD
```

### Step 4: Determine Order

Check the current order of source and target children in the LCA's children list.

```
Children order: [ResourceB, ResourceC, ResourceD]
Pattern: [*, SourceParent, *, TargetParent, *]
```

Where:
- `*` represents zero or more other children
- `SourceParent` = ResourceB (contains Source)
- `TargetParent` = ResourceD (contains Target)

This pattern `[*, SourceParent, *, TargetParent, *]` means **SourceParent is positioned to the left of TargetParent**.

Alternative pattern `[*, TargetParent, *, SourceParent, *]` would mean **TargetParent is positioned to the left of SourceParent**.

### Step 5: Calculate Movement Direction

Based on the order determined in Step 4, calculate which direction each resource should move to minimize distance:

**Case 1: Pattern `[*, SourceParent, *, TargetParent, *]` (Source is left of Target)**
- **Source side**: Link exits moving RIGHT (toward Target)
  - Move Source's ancestor to the RIGHTMOST position within its parent
- **Target side**: Link enters from LEFT (from Source)
  - Move Target's ancestor to the LEFTMOST position within its parent

**Case 2: Pattern `[*, TargetParent, *, SourceParent, *]` (Target is left of Source)**
- **Source side**: Link exits moving LEFT (toward Target)
  - Move Source's ancestor to the LEFTMOST position within its parent
- **Target side**: Link enters from RIGHT (from Source)
  - Move Target's ancestor to the RIGHTMOST position within its parent

### Step 6: Reorder Intermediate Resources

For each resource on the path from Source to LCA, and from Target to LCA, check if reordering is needed.

#### Case A: Same Direction as LCA + UnorderedChildren=true
Move the child to the edge that minimizes overlap with the link path.

**Example (horizontal layout, Source left of Target):**
```
ResourceB (contains Source=Resource1):
  - Direction: horizontal (same as LCA)
  - UnorderedChildren: true
  - Link exits moving RIGHT (toward Target)
  - Action: Move Resource1 to the RIGHTMOST position
  - Result: Link exits cleanly without overlapping Resource2

ResourceD (contains Target=Resource6):
  - Direction: horizontal (same as LCA)
  - UnorderedChildren: true
  - Link enters from LEFT (from Source)
  - Action: Move Resource6 to the LEFTMOST position
  - Result: Link enters cleanly without overlapping Resource5
```

**Why this works:**
- When the link exits Resource1 moving right, placing Resource1 at the rightmost position means no siblings (like Resource2) are in the link's path
- When the link enters Resource6 from the left, placing Resource6 at the leftmost position means no siblings (like Resource5) are in the link's path

#### Case B: Different Direction from LCA + UnorderedChildren=true
Move the child to the first position (index 0).

**Example (LCA is vertical, intermediate is horizontal):**

When the LCA has a different direction from its children, moving children to the first position helps reduce link crossing.

![Different Direction Before](./different-direction-before.png)

*Before: VPC (vertical) contains Horizontal1 and Horizontal2 (horizontal). Two links cross each other.*

![Different Direction After](./different-direction-after.png)

*After: Resource2 moved to first position in Horizontal1. Links no longer cross.*

**Hierarchy comparison:**

![Different Direction Tree Before](./different-direction-tree-before.png)

*Before: Original order VPC[H1,H2], H1[1,2], H2[3,4]. Links cross.*

![Different Direction Tree After](./different-direction-tree-after.png)

*After: Reordered H1[2,1], H2[3,4]. Green dashed line shows swap. Links no longer cross.*

**Why this works:**
- VPC direction: vertical (children stacked top-to-bottom)
- Horizontal1/2 direction: horizontal (different from LCA)
- By moving children to the first position, resources are consistently ordered
- This reduces the probability of link crossing in complex nested structures

### Step 7: Reorder at LCA Level

Finally, reorder children at the LCA level to place Target's parent adjacent to Source's parent. This minimizes the distance the link must travel and reduces overlap with intermediate resources.

**Pattern 1: Source before Target (Source is left of Target)**
```
Before: [*1, SourceParent(ResourceB), *2, TargetParent(ResourceD), *3]
After:  [*1, SourceParent(ResourceB), TargetParent(ResourceD), *2, *3]
```
Move TargetParent to the right of SourceParent, eliminating intermediate resources (*2) from the link path.

**Pattern 2: Target before Source (Target is left of Source)**
```
Before: [*1, TargetParent(ResourceD), *2, SourceParent(ResourceB), *3]
After:  [*1, *2, TargetParent(ResourceD), SourceParent(ResourceB), *3]
```
Move TargetParent to the left of SourceParent, eliminating intermediate resources (*2) from the link path.

**Key principle:** Keep SourceParent fixed, move TargetParent adjacent to it in the direction that maintains their relative order.

## Implementation Details

### Execution Timing

The reordering must occur after all links are registered but before layout calculation:

```
1. associateChildren()      // Add children to resources
2. associateLinks()          // Register all links
3. reorderChildrenByLinks()  // ← NEW: Reorder based on links
4. canvas.Scale()            // Calculate layout positions
```

### Processing Rules

- **Each link is processed exactly once** - No iteration needed
- **Only resources with UnorderedChildren=true are reordered**
- **BorderChildren are treated as if links originate from their parent**
- **Edge cases:**
  - 0 or 1 children: No reordering needed
  - No links: No reordering needed
  - Source and Target in same LCA child: No reordering needed

### Data Structure Operations

Children reordering operations:
- **Move to front:** O(1) with slice operations
- **Move to back:** O(1) with slice operations
- **Move adjacent:** O(n) where n is number of children

### Complexity

- **Time:** O(L × H) where L = number of links, H = tree height
- **Space:** O(H) for LCA ancestor tracking

## YAML Configuration

### Basic Usage

```yaml
Resources:
  ParentResource:
    Type: AWS::EC2::VPC
    Direction: horizontal
    Options:
      UnorderedChildren: true
    Children:
      - Child1
      - Child2
      - Child3
```

### Nested Example

```yaml
Resources:
  VPC:
    Type: AWS::EC2::VPC
    Direction: vertical
    Options:
      UnorderedChildren: true
    Children:
      - SubnetGroup1
      - SubnetGroup2
  
  SubnetGroup1:
    Type: AWS::Diagram::HorizontalStack
    Options:
      UnorderedChildren: true
    Children:
      - Subnet1
      - Subnet2
  
  SubnetGroup2:
    Type: AWS::Diagram::HorizontalStack
    Options:
      UnorderedChildren: true
    Children:
      - Subnet3
      - Subnet4

Links:
  - Source: Instance1  # In Subnet1
    Target: Instance4  # In Subnet4
```

## Benefits

1. **Automatic Optimization:** No manual reordering of YAML children needed
2. **Link Clarity:** Reduces visual overlap and crossing
3. **Maintainability:** Diagram structure adapts as links change
4. **Backward Compatible:** Only affects resources with UnorderedChildren=true

## Limitations

1. **Circular Links:** If links form cycles, the algorithm processes each link once without detecting cycles
2. **Conflicting Links:** Multiple links may have competing reordering preferences; processed in link order
3. **Manual Override:** Cannot manually specify order when UnorderedChildren=true

## Future Enhancements

- Cycle detection and warning
- Cost-based optimization for conflicting links
- Partial reordering (specify which children are movable)

## Related Features

- **LCA (Lowest Common Ancestor):** Used for finding common parent resources
- **Auto-positioning:** Automatic link position calculation
- **Orthogonal Links:** Right-angle link routing that benefits from optimal ordering

## References

- Issue #223: Orthogonal link improvement
- `internal/types/link.go`: LCA implementation
- `internal/types/resource.go`: Resource structure
