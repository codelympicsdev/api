(function() {
  var one = Math.random() * 10;
  var two = Math.random() * 10;
  return {
    Input: {
      Args: [],
      Stdin: one.toString() + '\n' + two.toString(),
    },
    Output: {
      Stdout: (one + two).toString(),
      Stderr: '',
    },
  };
})()
