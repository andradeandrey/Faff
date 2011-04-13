function configMathJax() {
        MathJax.Hub.Config({
                        extensions: ["tex2jax.js","TeX/AMSmath.js","TeX/AMSsymbols.js"],
                        jax: ["input/TeX", "output/HTML-CSS"],
                        tex2jax: {
			processEscapes: true,
                        inlineMath: [ ['$','$'], ["\\(","\\)"] ],
                        displayMath: [ ['$$','$$'], ["\\[","\\]"] ],
                },
                "HTML-CSS": { availableFonts: ["TeX"] }
        });
}
