**Reputation and trust mechanisms**

In this section, we will introduce a matrix based on a reputation matrix to guide agents in determining whether to trust other agents or not.
This is a particularly useful step since for a self-organized system to succeed: the way reputation is formed and dealt with by agents is indeed central to obtain a well-functioning self-organizing system since it provides a comprehensive metric of the social behavior of agents throughout the game, quickly allowing to identify selfish conduct and cooperative behaviour and their effects on the game and agents' welfare.

More importantly, the reputation mechanism truly provides an incentive for agents to cooperate, potentially maximizing their overall welfare.

## Reputation mechanism

### 1. Track Energy Usage per Turn

For each agent $i$, track the amount of energy they use in each turn $t$. Let $E_{i,t}$ denote the energy of the agent $i$ in the initial stage of turn $t$.

Define energy provision $p_{i,t}=E_{i,t}-E_{i,t-1}$

### 2. Account for Abandoning a Bike

Introduce a binary variable $A_{i,t}$ for each agent $i$ and turn $t$, which is 0 if the agent left a bike in that turn, and 1 otherwise.

### 3. Initialize the Reputation Matrix

Initialize the $N \times N$ matrix $M$ where $N$ is the total number of agents. Initially, all off-diagonal values are set to a mid-point (e.g., 0.5) and the diagonal values (self-reputation) to 1.

$$ M_{initial} = \begin{bmatrix}
1 & 0.5 & ... & 0.5 \\
0.5 & 1 & ... & 0.5 \\
... & ... & ... & ... \\
0.5 & 0.5 & ... & 1
\end{bmatrix} $$

### 4. Calculate a Base Reputation Score
Define Utility:

$$ U_{i,t} = \frac{a}{n} \sum_{j=1}^n p_j + b(E_{i,t-1} - p_i)$$ 

Overall utility function conception 
Where a>b and $\frac{a}{n} < b$. a,b are both parameters, p and E defined previously. 
If agent is a freerider, $p_i=0$, this might be a large value, vv. It needs to be **normlized** into 0-1 values.


Parameters and where to get input values
Citation to the slide(s) where the concept is explained
(Plus, if you can) How to merge with the Reputation vector to make it adaptive for each agent

Define Benefit rate:

$$ B_{i,t}=p_k-\frac{1}{n-1}\sum_{k=1,~k\neq i}^n p_k $$ 

$B_{i,t}$ negative value when agent i is a freerider,vv.

The base reputation score $R_{i,t}$ now considers both utilarity and the act of abandoning a bike. It can be calculated as:

$$ R_{i,t} = U_{i,t} + \beta A_{i,t} $$

Here, $\beta$ is the weight given to the act of abandoning a bike, reflecting its impact on the reputation. It needs to be **normlized** into 0-1 values.


### 5. Update the Matrix After Each Turn

The update rule now also considers $A_{i,t}$. Update the matrix based on the new energy usage data and the act of leaving a bike. The update rule for element $M_{i,j}$ in the matrix can be:

$$
M_{i,j,t+1} =
\begin{cases}
\delta M_{i,j,t} \cdot + (1 - \delta)(B_{i,t}\times(1-R_{j,t}))& \text{if } i \neq j \\
1 & \text{if } i = j
\end{cases}
$$

Here, $\alpha$ is a factor determining how much the reputation changes per turn.


### 6. Normalization and Adjustment

Normalize the matrix after each update to ensure that the reputation scores are comparable across agents and turns:

$$ M_{i,j,t+1} = \frac{M_{i,j,t+1}}{\sum_{k=1}^{N} M_{i,k,t+1}} $$

*Considerations:*
- This model assumes a direct correlation between energy usage and positive contribution, which might need adjustments based on game dynamics.
- The model can be further refined by incorporating additional metrics or observations from the game.
- The value of $\beta$ should be chosen to reflect how significantly abandoning a bike is perceived compared to energy contribution.
