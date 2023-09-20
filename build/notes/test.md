<!DOCTYPE html>
<html>

<head>
    <title>Ian Kilty's Website</title>
    <link rel="preconnect" href="https://fonts.googleapis.com">
    <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
    <link
        href="https://fonts.googleapis.com/css2?family=Source+Serif+4:ital,opsz,wght@0,8..60,300;0,8..60,400;0,8..60,500;0,8..60,600;1,8..60,300;1,8..60,400;1,8..60,500;1,8..60,600&family=VT323&display=swap"
        rel="stylesheet">
    
    
    <style>
        html {
            font-family: 'Source Serif 4', serif;
        }

        .black {
            color: #000;
        }

        .contact-links {
            border-bottom: solid 1px #F2F2F2;
            border-bottom-width: 1px;
            border-bottom-style: solid;
            border-bottom-color: rgb(242, 242, 242);
        }

        .links {
            font-size: 18px;
            display: flex;
            color: black;
        }

        .links a {
            padding-right: 10px;
            color: black;
        }

        footer {
            padding-top: 15px;
            font-size: 17px;
            justify-content: space-around;
            display: flex;
            padding-bottom: 120px;
        }

        a {
            color: black !important;
        }

        body {
            background: #fff;
        }
    </style>
</head>

<body id="main">
    <main style="width:50%; margin:auto;">
        <script>
            var color = true;
            var toggle_color = function () {
                const main = document.getElementById("main")
                if (color) {
                    main.style.background = "#fff"
                    main.style.color = "#000"
                } else {
                    main.style.background = "#1e2021"
                    main.style.color = "rgb(193, 193, 193)"
                }
                color = !color
            }
        </script>
        <h2 onclick="toggle_color()">Ian Kilty's Website</h2>
        <div class="links black" style="font-size: 17px; text-decoration: none;">
            <a href="/">Home</a>
            <a href="/blog">Blog</a>
            <a href="/projects">Projects</a>
            <a href="/notes">Notes</a>
        </div>
        <h1>testing</h1>
        <p><em>9/20/2023</em></p>
        <p>
            Tags:
            <a href="">Tags</a>
            <a href="">For</a>
            <a href="">The</a>
            <a href="">Article</a>
        </p>
        <article
            style="font-size: 20px; line-height: 32px; text-rendering: optimizeLegibility; letter-spacing: -0.06px; margin-bottom: -9.2px;">
            
        </article>
        <br>
        <hr>
        <footer>
            <a href="">Back to the top</a>
            <a href="">Github</a>
            <a href="">Linkedin</a>
            <a href="">Email</a>
        </footer>
    </main>
</body>

</html>