"use strict";

const e = React.createElement;

class LikeButton extends React.Component {

  constructor(props) {
    super(props);
    this.state = { liked: false };
  }

  render() {
    if (this.state.liked) {
      return "You liked this.";
    }

    return e(
      "button",
      { onClick: () => this.setState({ liked: true }) },
      "Like"
    );
  }
}

function LikeButton2(){
  const [v, setV] = React.useState(34)
  return <button>
    likeme {v}
  </button>
}

const domContainer = document.querySelector("#react_container");
const root = ReactDOM.createRoot(domContainer);
root.render(e(LikeButton2));
