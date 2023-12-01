## Reputation mechanism

### 1. Track Energy Usage per Turn

For each agent $i$, record the energy they use in each turn $t$, denoted as $E_{i,t}$.

### 2. Account for Abandoning a Bike

For each agent $i$ in each turn $t$, have a binary indicator $A_{i,t}$ that is 1 if the agent left a bike in that turn, and 0 otherwise.

### 3. Measure Distance to Loot Boxes

For each agent $i$ in a leadership position at turn $t$, measure the distance from their chosen loot box to the closest loot box, denoted as $D_{i,t}$. This will reflect the leader's decision-making with respect to team benefit versus personal gain.

### 4. Calculate a Base Reputation Score

Calculate the base reputation score $R_{i,t}$ considering energy use, bike abandonment, and decision-making regarding loot box choice. The score can be calculated as:

$$ R_{i,t} = \left(1 - \frac{E_{i,t}}{\max(E_{1,t}, E_{2,t}, ..., E_{N,t})}\right) \times \left(1 - \beta \times A_{i,t}\right) \times \left(1 - \gamma \times \frac{D_{i,t}}{D_{\text{closest},t}}\right) $$

Here, $\beta$ is the weight given to the act of abandoning a bike, $\gamma$ is the weight given to the penalty for directing towards a suboptimal loot box, and $D_{\text{closest},t}$ is the distance to the closest loot box.

### 5. Initialize the Reputation Vector

Initialize the reputation vector $R$ with a length of $N$, where $N$ is the total number of agents. Initially, all values can be set to a neutral value, such as 0.5.

$$ R_{initial} = [0.5, 0.5, ..., 0.5] $$

### 6. Update the Vector After Each Turn

After each turn, update the vector based on the new data. The update for agent $i$'s reputation can be:

$$ R_{i,t+1} = R_{i,t} \cdot \left(1 + \alpha \cdot (R_{i,t} - 0.5)\right) - \lambda \cdot A_{i,t} - \xi \cdot \frac{D_{i,t}}{D_{\text{closest},t}} $$

Here, $\alpha$ is a factor determining how much the reputation changes due to energy consumption, $\lambda$ is the penalty for abandoning a bike, and $\xi$ is the penalty for directing the megabike towards a suboptimal loot box.

### 7. Apply Decay Factor and Normalize

Apply a decay factor $\delta$ to account for the diminishing influence of past actions:

$$ R_{i,t+1} = \delta \cdot R_{i,t+1} + (1 - \delta) \cdot R_{i,t} $$

Normalize the reputation vector so that the sum of all reputations equals 1:

$$ R_{i,t+1} = \frac{R_{i,t+1}}{\sum_{j=1}^{N} R_{j,t+1}} $$

The resulting reputation vector $R_{t+1}$ provides a comprehensive score for each agent, taking into account their energy contribution, their loyalty to the team in terms of bike abandonment, and their leadership in directing the megabike towards loot boxes in a way that reflects collective benefit.

*Considerations:*
- This model assumes a direct correlation between energy usage and positive contribution, which might need adjustments based on game dynamics.
- The model can be further refined by incorporating additional metrics or observations from the game.
- The value of $\beta$ should be chosen to reflect how significantly abandoning a bike is perceived compared to energy contribution.

## Decision Making

The produced reputation matrix would be used at each turn to weight important decisions such as:
- Removing an agent from the mega-bike (e.g., if overall reputation from the other agents towards the potentially removable one is below a threshold)
- Allocating energy (proportional to aggregated reputation value of each agent)
- Lootbox decision making (e.g., Leader being the agent with the highest aggregated reputation value)
- Accepting an agent to the mega-bike (select the best if multiple applicants, accept above threshold if single applicant)

## Experiments

The reputation matrix introduces several parameters that should be tuned (e.g., pedalling contribution weight, bike-abandonment weight), and others that have to be derived according to which decisions reputation should impact (e.g., threshold for agent removal/acceptance, leader voting technique). This way, several experiments could be performed to analyze various social scenarios.
