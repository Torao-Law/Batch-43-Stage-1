// let namaVariable = "isivariable"

// const angka = 5

// const isMarried = true

// type data non primitive
// Array
// let array = ["Dandi", 17, false]

// console.log(array)
// console.table(array[0])

// object
// let myBio = {
//     name: "Dandi Saputra",
//     age: 17,
//     isMarried: {
//         age2: 24
//     }
// }

// console.table(`Apakah mas dandi sudah menikah ? ${myBio.isMarried.age2}`)

// // array of object
// let myFriends = [
//     {
//         name: "Nima",
//         age: 23,
//         address: "Tangerang",
//         isMarried: true
//     },
//     {
//         name: "Agung",
//         age: 32,
//         address: "Jakarta",
//         isMarried: true,
//         istri: "angle"
//     },
//     {
//         name: "Krisna",
//         age: 15,
//         address: "Jakarta",
//         isMarried: true,
//         istri: "merri",
//         myCar: "supra"
//     }
// ]

// console.table(myFriends[0].myCar)

// console.table(myFriends[0].address)

// let example = document.getElementById("example")
// example.innerHTML += `
// <div>
//     <img src=${image} />
// </div>
// `


// 1. nilai dasar/awal
// 2. kondisi
// 3. initstatemen dan post statement

// for(let i = 1; i<=10; i++) {
//     console.table(`Hello world ${i}`)
// }

let blogs = []

function getData(event) {
    event.preventDefault()

    let title = document.getElementById("input-blog-title").value
    let content = document.getElementById("input-blog-content").value
    let image = document.getElementById("input-blog-image").files

    image = URL.createObjectURL(image[0])

    let addBlog = {
        title,
        content,
        image,
        author: "Dandi Saputra",
        postedAt: new Date()
    }

    blogs.push(addBlog)

    console.log(blogs)
    showData()
}

function showData() {
    document.getElementById("contents").innerHTML = ""

    for(let i= 0; i <= blogs.length; i++){
        document.getElementById("contents").innerHTML += `
        <div class="blog-list-item">
        <div class="blog-image">
          <img src="${blogs[i].image}" alt="" />
        </div>
        <div class="blog-content">
          <div class="btn-group">
            <button class="btn-edit">Edit Post</button>
            <button class="btn-post">Post Blog</button>
          </div>
          <h1>
            <a href="blog-detail.html" target="_blank"
              >${blogs[i].title}</a
            >
          </h1>
          <div class="detail-blog-content">
            ${blogs[i].postedAt} | ${blogs[i].author}
          </div>
          <p>${blogs[i].content}
          </p>

          <div style="float:right; margin: 10px">
            <p style="font-size: 15px; color:grey">1 minutes ago</p>
          </div>

        </div>
      </div>
        `
    }
}



