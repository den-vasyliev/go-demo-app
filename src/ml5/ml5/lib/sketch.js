// Copyright (c) 2018 ml5
// 
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

/* ===
ml5 Example
Style Transfer Image Example
This uses a pre-trained model of Udnie (Young American Girl, The Dance)
=== */
const inputImg = document.getElementById('inputImg'); // The image we want to transfer
const statusMsg = document.getElementById('statusMsg'); // The status message
const Udnie = document.getElementById('Udnie'); // The div contrianer that holds new style image A

ml5.styleTransfer('models/udnie')
  .then(style => style.transfer(inputImg))
  .then(result => {
    const newImage = new Image(250, 250);
    newImage.src = result.src;
    Udnie.appendChild(newImage);
    statusMsg.innerHTML = 'Done!';
  });
