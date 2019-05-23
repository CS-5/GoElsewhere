import React from "react";
import Button from "@material-ui/core/Button";
import Dialog from "@material-ui/core/Dialog";
import DialogActions from "@material-ui/core/DialogActions";
import DialogContent from "@material-ui/core/DialogContent";
import DialogContentText from "@material-ui/core/DialogContentText";
import DialogTitle from "@material-ui/core/DialogTitle";
import withMobileDialog from "@material-ui/core/withMobileDialog";
import TextField from "@material-ui/core/TextField";

function CreateDialog({ fullScreen, open, closer, creater }) {
  function handleSubmit(event) {
    event.preventDefault();
    const data = new FormData(event.target);

    creater(data);
  }

  return (
    <Dialog
      fullScreen={fullScreen}
      open={open}
      onClose={closer}
      aria-labelledby="create-redirect"
    >
      <DialogTitle id="create-redirect">{"Create a new Redirect"}</DialogTitle>
      <DialogContent>
        <DialogContentText>
          Enter the desired URL and code (Optional)
        </DialogContentText>
        <form
          id="createForm"
          onSubmit={e => {
            handleSubmit(e);
            closer();
          }}
        >
          <TextField
            id="outlined-url"
            label="URL"
            name="url"
            margin="normal"
            fullWidth={true}
            variant="outlined"
            type="url"
            required
          />
          <TextField
            id="outlined-code"
            label="Code"
            name="code"
            margin="normal"
            fullWidth={true}
            variant="outlined"
            type="text"
          />
        </form>
      </DialogContent>
      <DialogActions>
        <Button onClick={closer} color="secondary">
          Cancel
        </Button>
        <Button type="submit" form="createForm" color="secondary" autoFocus>
          Create
        </Button>
      </DialogActions>
    </Dialog>
  );
}

export default withMobileDialog()(CreateDialog);
