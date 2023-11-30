**Reputation and trust mechanisms**

In this section, we will introduce a matrix based on a reputation matrix to guide agents in determining whether to trust other agents or not.
This is a particularly useful step since for a self-organized system to succeed: the way reputation is formed and dealt with by agents is indeed central to obtain a well-functioning self-organizing system since it provides a comprehensive metric of the social behavior of agents throughout the game, quickly allowing to identify selfish conduct and cooperative behaviour and their effects on the game and agents' welfare.

More importantly, the reputation mechanism truly provides an incentive for agents to cooperate, potentially maximizing their overall welfare.

# Mauro:

## Reputation mechanism

### 1. Track Energy Usage per Turn

For each agent $i$, track the amount of energy they use in each turn $t$. Let $E_{i,t}$ denote the energy used by agent $i$ in turn $t$.

### 2. Account for Abandoning a Bike

Introduce a binary variable $A_{i,t}$ for each agent $i$ and turn $t$, which is 1 if the agent left a bike in that turn, and 0 otherwise.

### 3. Calculate a Base Reputation Score

The base reputation score $R_{i,t}$ now considers both energy use and the act of abandoning a bike. It can be calculated as:

$$ R_{i,t} = (1 - \frac{E_{i,t}}{\max(E_{1,t}, E_{2,t}, ..., E_{N,t})}) \times (1 - \beta \times A_{i,t}) $$

Here, $\beta$ is the weight given to the act of abandoning a bike, reflecting its impact on the reputation.

### 4. Initialize the Reputation Matrix

Initialize the $N \times N$ matrix $M$ where $N$ is the total number of agents. Initially, all off-diagonal values are set to a mid-point (e.g., 0.5) and the diagonal values (self-reputation) to 1.

$$ M_{initial} = \begin{bmatrix}
1 & 0.5 & ... & 0.5 \\
0.5 & 1 & ... & 0.5 \\
... & ... & ... & ... \\
0.5 & 0.5 & ... & 1
\end{bmatrix} $$

### 5. Update the Matrix After Each Turn

The update rule now also considers $A_{i,t}$. Update the matrix based on the new energy usage data and the act of leaving a bike. The update rule for element $M_{i,j}$ in the matrix can be:

$$
M_{i,j,t+1} =
\begin{cases}
M_{i,j,t} \cdot (1 + \alpha \cdot (R_{j,t} - 0.5)) & \text{if } i \neq j \\
1 & \text{if } i = j
\end{cases}
$$

Here, $\alpha$ is a factor determining how much the reputation changes per turn.

### 6. Incorporate a Decay Factor

Apply the decay factor $\delta$ to ensure that the impact of an agent's actions on their reputation diminishes over time:

$$ M_{i,j,t+1} = \delta \cdot M_{i,j,t+1} + (1 - \delta) \cdot M_{i,j,t} $$

### 7. Normalization and Adjustment

Normalize the matrix after each update to ensure that the reputation scores are comparable across agents and turns:

$$ M_{i,j,t+1} = \frac{M_{i,j,t+1}}{\sum_{k=1}^{N} M_{i,k,t+1}} $$

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

# Yanzhou:

**Calculation of Reputation Matrix**

In the system, there will be n agents. For an agent named 'Monkey,' its reputation matrix is considered as one of its attributes. This can be denoted as $R_{\text{Monkey}}$ or simply referred to as $R_{m}$.

$$R_{m}=[r1, r2, ..., ri=1, ..., rn]$$
j
**Iterations**

We refer to a complete cycle as one iteration or round loop. For the i th iteration, it is denoted as iter i.


## Reputation Update Rule

The reputation update from agent \(i\) to agent \(j\) can be expressed using the following formula:

```math
R_m(j)_{\text{k}} = \alpha \cdot R_m(j)_{\text{k-1}} + \beta \cdot \text{Feedback}_i(j)
Feedback=Utility
```

where $\alpha$ is the decay factor, used to diminish the impact of past reputations, and $\beta$ is the learning rate, adjusting the influence of new feedback on reputation. The function Feedback(i, j) represents feedback obtained from the environment or other agents regarding the interaction between Monkey $m$ and $j$.

## Trust and Distrust Update Rules

Trust and distrust can be updated with the following formulas:
```math
T(i, j)_{\text{new}} = \gamma \cdot T(i, j)_{\text{old}} + \delta \cdot \text{GoodBehavior}(i, j)

D(i, j)_{\text{new}} = \rho \cdot D(i, j)_{\text{old}} + \sigma \cdot \text{BadBehavior}(i, j)
```
Here, \(\gamma\) and \(\rho\) are decay factors, and \(\delta\) and \(\sigma\) are learning rates. GoodBehavior(i, j) and BadBehavior(i, j) represent metrics of positive and negative behaviors of agent \(i\) towards agent \(j\).

## Cooperation Decision Rule

When considering cooperation with agent \(j\), agent \(i\) can use the following rule:
```math
\text{Cooperate}(i, j) = \frac{T(i, j)}{T(i, j) + D(i, j)}
```

This value ranges from 0 to 1, indicating the likelihood of cooperation. A higher trust level results in a value closer to 1, indicating a higher likelihood of cooperation.
