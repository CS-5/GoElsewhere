import React from "react";
import Grid from "@material-ui/core/Grid";
import { Fab, Box, CircularProgress, Tooltip } from "@material-ui/core";
import AddIcon from "@material-ui/icons/Add";
import RedirectCard from "./RedirectCard";

export default function RedirectList({
  redirects = [],
  createDialog,
  deleteDialog
}) {
  return (
    <>
      {redirects.length === 0 ? (
        <Box
          width={1}
          display="flex"
          justifyContent="center"
          alignItems="center"
        >
          <CircularProgress />
        </Box>
      ) : (
        <Grid container spacing={2}>
          {redirects.map(redirect => {
            return (
              <Grid key={redirect.code} item xs={12} sm={6} md={3}>
                <RedirectCard
                  key={redirect.id}
                  code={redirect.code}
                  url={redirect.url}
                  created={redirect.created}
                  hits={redirect.hits}
                  link={redirect.link}
                  deleter={deleteDialog}
                />
              </Grid>
            );
          })}
        </Grid>
      )}
      <Box position="fixed" bottom={16} right={16}>
        <Tooltip title="Create" placement="left">
          <Fab onClick={createDialog} color="secondary">
            <AddIcon />
          </Fab>
        </Tooltip>
      </Box>
    </>
  );
}
