The derivative of a constant is always zero and the derivative of x is always 1.

$$ \frac{d}{dx} c = 0 $$
$$ \frac{d}{dx} x = 1 $$

Power rule

$$ \frac{d}{dx} x^n = nx^{n-1} $$
$$ \frac{d}{dx} a^x = a^x \cdot x' \cdot ln(a) $$
But the derivative of x is 1 so we write
$$ \frac{d}{dx} a^x = a^x \cdot ln(a) $$
Where u is a function of x
$$ \frac{d}{dx} a^u = a^u \cdot u' \cdot ln(a) $$
$$ \frac{d}{dx} cf(x) = c \cdot f'(x) $$
$$ \frac{d}{dx} 5x^4 = 5 \cdot \frac{d}{dx} x^4 = 5 \cdot 4x^3 = 20x^3 $$
Product rule
$$ \frac{d}{dx} (uv) = u'v + uv' $$
Quotient rule
$$ \frac{d}{dx} \frac u v = \frac {u'v - uv'}{v^2} $$
Chain rule
$$ \frac{d}{dx} f(g(x)) = f'(g(x)) \cdot g'(x) \cdot x' = f'(g(x)) \cdot g'(x) $$
Chain rule with power rule
$$ \frac{d}{dx} f(x)^n = nf(x)^{n-1} \cdot f'(x) $$
Logarithmic functions
$$ \frac{d}{dx} log_a(u) = \frac {u'}{u\ln(a)} $$
$$ \frac{d}{dx} \ln(u) = \frac {u'}{u\ln(e)} = \frac {u'}{u} $$
Trigonometric functions
$$ \frac {d}{dx}\sin(x) = \cos(x)x' = \cos(x)$$
$$ \frac {d}{dx}\sin^{-1}(x) = \frac{x'}{\sqrt{1-x^2}}$$
$$ \frac {d}{dx}\cos(x) = -\sin(x)x' $$
$$ \frac {d}{dx}\cos^{-1}(x) = \frac{-x'}{\sqrt{1-x^2}}$$
$$ \frac {d}{dx}\tan(x) = \sec^2(x)x'$$
$$ \frac {d}{dx}\tan^{-1}(x) = \frac{x'}{1+x^2}$$
$$ \frac {d}{dx}\cot(x) = -\csc^2(x)x'$$
$$ \frac {d}{dx}\cot^{-1}(x) = \frac{-x'}{1+x^2}$$
$$ \frac {d}{dx}\sec(x) = \sec(x)\tan(x)x'$$
$$ \frac {d}{dx}\sec^{-1}(x) = \frac{x'}{x\sqrt{x^2-1}}$$
$$ \frac {d}{dx}\csc(x) = -\csc(x)\cot(x)x' $$
$$ \frac {d}{dx}\csc^{-1}(x) = \frac{-x'}{x\sqrt{x^2-1}}$$
You need to use the chain rule if you have something besides x.
$$ \frac {d}{dx}\sin(x^2) = \cos(x^2)2x $$

#### Logarithmic differentiation
How can we solve for the derivative of a variable raised to a variable?
$$ \frac{d}{dx} x^x $$
$$ y = x^x $$
Before you can take the derivative of both sides, you need to take the natural logarithm to move the exponent down.
$$ \ln(y) = \ln(x^x) = x\ln(x) $$
$$ \ln(y) = x\ln(x) $$
$$ \frac{dy}{dx} \ln(y) = \frac{d}{dx} x\ln(x) $$
$$ \frac{1}{y} \cdot \frac{dy}{dx} = (1)\ln(x) + x(\frac{1}{x}) = \ln(x) + 1 $$
$$ y\cdot(\frac{1}{y} \cdot \frac{dy}{dx}) = (\ln(x) + 1)\cdot y $$
$$ \frac{dy}{dx} = y(\ln(x) + 1) $$
$$ \frac{dy}{dx} = x^x(\ln(x) + 1) $$




